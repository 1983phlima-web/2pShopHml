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

const selectCols = `id, tenant_id, email, name, role, password_hash, avatar, phone, active, created_at, updated_at`

func scanUser(row pgx.Row, u *domain.User) error {
	var avatar, phone *string
	err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.Name, &u.Role, &u.Password, &avatar, &phone, &u.Active, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return err
	}
	if avatar != nil {
		u.Avatar = *avatar
	}
	if phone != nil {
		u.Phone = *phone
	}
	return nil
}

func (r *Repository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO users (id, tenant_id, email, name, role, password_hash, avatar, phone, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULLIF($8, ''), $9, $10, $11)
	`, user.ID, user.TenantID, user.Email, user.Name, user.Role, user.Password, user.Avatar, user.Phone, user.Active, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *Repository) GetByID(ctx context.Context, tenantID, id string) (*domain.User, error) {
	var u domain.User
	row := r.db.Pool.QueryRow(ctx, `SELECT `+selectCols+` FROM users WHERE tenant_id = $1 AND id = $2`, tenantID, id)
	if err := scanUser(row, &u); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New(errors.ErrNotFound)
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetByEmail(ctx context.Context, tenantID, email string) (*domain.User, error) {
	var u domain.User
	row := r.db.Pool.QueryRow(ctx, `SELECT `+selectCols+` FROM users WHERE tenant_id = $1 AND email = $2`, tenantID, email)
	if err := scanUser(row, &u); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New(errors.ErrNotFound)
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) ListByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*domain.User, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT `+selectCols+` FROM users WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := scanUser(rows, &u); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, rows.Err()
}

func (r *Repository) UpdateAvatar(ctx context.Context, tenantID, userID, avatar string) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE users SET avatar = $3, updated_at = NOW() WHERE tenant_id = $1 AND id = $2`, tenantID, userID, avatar)
	return err
}
