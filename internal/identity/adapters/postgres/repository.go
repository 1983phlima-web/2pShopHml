package postgres

import (
	"context"

	"github.com/2pshop/2pshop/internal/identity/domain"
	"github.com/2pshop/2pshop/internal/identity/ports"
	"github.com/2pshop/2pshop/internal/platform"
	"github.com/2pshop/2pshop/pkg/errors"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *platform.DB
}

func NewRepository(db *platform.DB) ports.Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO users (id, tenant_id, email, name, role, password_hash, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, user.ID, user.TenantID, user.Email, user.Name, user.Role, user.Password, user.Active, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *Repository) GetByID(ctx context.Context, tenantID, id string) (*domain.User, error) {
	var u domain.User
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, tenant_id, email, name, role, password_hash, active, created_at, updated_at
		FROM users WHERE tenant_id = $1 AND id = $2
	`, tenantID, id).Scan(&u.ID, &u.TenantID, &u.Email, &u.Name, &u.Role, &u.Password, &u.Active, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetByEmail(ctx context.Context, tenantID, email string) (*domain.User, error) {
	var u domain.User
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, tenant_id, email, name, role, password_hash, active, created_at, updated_at
		FROM users WHERE tenant_id = $1 AND email = $2
	`, tenantID, email).Scan(&u.ID, &u.TenantID, &u.Email, &u.Name, &u.Role, &u.Password, &u.Active, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) ListByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*domain.User, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, tenant_id, email, name, role, password_hash, active, created_at, updated_at
		FROM users WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.TenantID, &u.Email, &u.Name, &u.Role, &u.Password, &u.Active, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, rows.Err()
}
