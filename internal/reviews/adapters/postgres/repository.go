package postgres

import (
	"context"

	"github.com/2pshop/2pshop/internal/platform"
	"github.com/2pshop/2pshop/internal/reviews/domain"
	"github.com/2pshop/2pshop/internal/reviews/ports"
)

type Repository struct {
	db *platform.DB
}

func NewRepository(db *platform.DB) ports.Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, rv *domain.Review) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO product_reviews (id, tenant_id, product_id, user_id, user_name, rating, comment, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, rv.ID, rv.TenantID, rv.ProductID, rv.UserID, rv.UserName, rv.Rating, rv.Comment, rv.CreatedAt)
	return err
}

func (r *Repository) ListByProduct(ctx context.Context, tenantID, productID string, limit, offset int) ([]*domain.Review, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, tenant_id, product_id, user_id, user_name, rating, comment, created_at
		FROM product_reviews WHERE tenant_id=$1 AND product_id=$2
		ORDER BY created_at DESC LIMIT $3 OFFSET $4
	`, tenantID, productID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*domain.Review
	for rows.Next() {
		var rv domain.Review
		if err := rows.Scan(&rv.ID, &rv.TenantID, &rv.ProductID, &rv.UserID, &rv.UserName, &rv.Rating, &rv.Comment, &rv.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, &rv)
	}
	return reviews, rows.Err()
}

// ListByUser returns every review the given user has written — the
// "Testemunhos" section of "Meus Comentários" — including the product
// name via a join, so the frontend doesn't need extra fetches.
func (r *Repository) ListByUser(ctx context.Context, tenantID, userID string, limit, offset int) ([]*domain.Review, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT rv.id, rv.tenant_id, rv.product_id, p.name, rv.user_id, rv.user_name, rv.rating, rv.comment, rv.created_at
		FROM product_reviews rv
		JOIN products p ON p.id = rv.product_id AND p.tenant_id = rv.tenant_id
		WHERE rv.tenant_id=$1 AND rv.user_id=$2
		ORDER BY rv.created_at DESC LIMIT $3 OFFSET $4
	`, tenantID, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*domain.Review
	for rows.Next() {
		var rv domain.Review
		if err := rows.Scan(&rv.ID, &rv.TenantID, &rv.ProductID, &rv.ProductName, &rv.UserID, &rv.UserName, &rv.Rating, &rv.Comment, &rv.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, &rv)
	}
	return reviews, rows.Err()
}

func (r *Repository) Summary(ctx context.Context, tenantID, productID string) (domain.Summary, error) {
	var avg *float64
	var count int
	err := r.db.Pool.QueryRow(ctx, `
		SELECT AVG(rating), COUNT(*) FROM product_reviews WHERE tenant_id=$1 AND product_id=$2
	`, tenantID, productID).Scan(&avg, &count)
	if err != nil {
		return domain.Summary{}, err
	}
	s := domain.Summary{ProductID: productID, Count: count}
	if avg != nil {
		s.Average = *avg
	}
	return s, nil
}
