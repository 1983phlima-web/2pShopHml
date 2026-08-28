package idempotency

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/2pshop/2pshop/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusCompleted Status = "COMPLETED"
	StatusFailed    Status = "FAILED"
)

type Key struct {
	TenantID     string    `json:"tenant_id"`
	Key          string    `json:"key"`
	Operation    string    `json:"operation"`
	RequestHash  string    `json:"request_hash"`
	Status       Status    `json:"status"`
	ResponseCode int       `json:"response_code,omitempty"`
	ResponseBody []byte    `json:"response_body,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Store interface {
	Get(ctx context.Context, tenantID, key string) (*Key, error)
	Save(ctx context.Context, k *Key) error
	Update(ctx context.Context, k *Key) error
}

type RedisStore struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

func NewRedisStore(client *redis.Client, ttl time.Duration) *RedisStore {
	return &RedisStore{
		client: client,
		prefix: "idempotency:",
		ttl:    ttl,
	}
}

func (s *RedisStore) redisKey(tenantID, key string) string {
	return fmt.Sprintf("%s%s:%s", s.prefix, tenantID, key)
}

func (s *RedisStore) Get(ctx context.Context, tenantID, key string) (*Key, error) {
	data, err := s.client.Get(ctx, s.redisKey(tenantID, key)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var k Key
	if err := json.Unmarshal(data, &k); err != nil {
		return nil, err
	}
	return &k, nil
}

func (s *RedisStore) Save(ctx context.Context, k *Key) error {
	data, err := json.Marshal(k)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, s.redisKey(k.TenantID, k.Key), data, s.ttl).Err()
}

func (s *RedisStore) Update(ctx context.Context, k *Key) error {
	return s.Save(ctx, k)
}

func HashPayload(payload any) (string, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:]), nil
}

type Manager struct {
	store Store
}

func NewManager(store Store) *Manager {
	return &Manager{store: store}
}

func (m *Manager) Check(ctx context.Context, tenantID, key, operation string, payload any) (*Key, error) {
	existing, err := m.store.Get(ctx, tenantID, key)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "idempotency store error", err)
	}

	hash, err := HashPayload(payload)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to hash payload", err)
	}

	if existing != nil {
		if existing.RequestHash != hash {
			return nil, errors.New(errors.ErrIdempotency).WithDetail("reason", "same key with different payload")
		}
		if existing.Status == StatusPending {
			return nil, errors.New(errors.ErrConflict).WithDetail("reason", "request still in progress")
		}
		return existing, nil
	}

	k := &Key{
		TenantID:    tenantID,
		Key:         key,
		Operation:   operation,
		RequestHash: hash,
		Status:      StatusPending,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		CreatedAt:   time.Now(),
	}
	if err := m.store.Save(ctx, k); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to save idempotency key", err)
	}
	return nil, nil
}

func (m *Manager) Complete(ctx context.Context, tenantID, key string, code int, body []byte) error {
	k, err := m.store.Get(ctx, tenantID, key)
	if err != nil {
		return err
	}
	if k == nil {
		return fmt.Errorf("idempotency key not found")
	}
	k.Status = StatusCompleted
	k.ResponseCode = code
	k.ResponseBody = body
	return m.store.Update(ctx, k)
}

func (m *Manager) Fail(ctx context.Context, tenantID, key string) error {
	k, err := m.store.Get(ctx, tenantID, key)
	if err != nil {
		return err
	}
	if k == nil {
		return fmt.Errorf("idempotency key not found")
	}
	k.Status = StatusFailed
	return m.store.Update(ctx, k)
}
