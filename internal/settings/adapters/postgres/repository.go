package postgres

import (
	"context"
	"encoding/json"

	"github.com/2pshop/2pshop/internal/platform"
	"github.com/2pshop/2pshop/internal/settings/domain"
	"github.com/2pshop/2pshop/internal/settings/ports"
	"github.com/2pshop/2pshop/pkg/errors"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *platform.DB
}

func NewRepository(db *platform.DB) ports.Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(ctx context.Context, tenantID, key string) (*domain.Setting, error) {
	var raw []byte
	err := r.db.Pool.QueryRow(ctx, `
		SELECT value FROM system_settings WHERE tenant_id=$1 AND key=$2
	`, tenantID, key).Scan(&raw)
	if err == pgx.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	var value map[string]any
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &value)
	}
	return &domain.Setting{Key: key, Value: value}, nil
}

func (r *Repository) Set(ctx context.Context, tenantID, key string, value map[string]any) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = r.db.Pool.Exec(ctx, `
		INSERT INTO system_settings (id, tenant_id, key, value, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, NOW())
		ON CONFLICT (tenant_id, key) DO UPDATE SET value = $3, updated_at = NOW()
	`, tenantID, key, raw)
	return err
}
