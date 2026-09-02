package postgres

import (
	"context"

	"github.com/2pshop/2pshop/internal/platform"
	"github.com/2pshop/2pshop/internal/questions/domain"
	"github.com/2pshop/2pshop/internal/questions/ports"
)

type Repository struct {
	db *platform.DB
}

func NewRepository(db *platform.DB) ports.Repository {
	return &Repository{db: db}
}

const selectCols = `id, tenant_id, product_id, user_id, user_name, question, answer, answered_at, created_at`

func (r *Repository) Create(ctx context.Context, q *domain.Question) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO product_questions (id, tenant_id, product_id, user_id, user_name, question, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, q.ID, q.TenantID, q.ProductID, q.UserID, q.UserName, q.Question, q.CreatedAt)
	return err
}

func scanQuestion(rows interface {
	Scan(dest ...any) error
}) (*domain.Question, error) {
	var q domain.Question
	if err := rows.Scan(&q.ID, &q.TenantID, &q.ProductID, &q.UserID, &q.UserName, &q.Question, &q.Answer, &q.AnsweredAt, &q.CreatedAt); err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *Repository) ListByProduct(ctx context.Context, tenantID, productID string, limit, offset int) ([]*domain.Question, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT `+selectCols+` FROM product_questions WHERE tenant_id=$1 AND product_id=$2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`, tenantID, productID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*domain.Question
	for rows.Next() {
		q, err := scanQuestion(rows)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, rows.Err()
}

func (r *Repository) ListByUser(ctx context.Context, tenantID, userID string, limit, offset int) ([]*domain.Question, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT `+selectCols+` FROM product_questions WHERE tenant_id=$1 AND user_id=$2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`, tenantID, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*domain.Question
	for rows.Next() {
		q, err := scanQuestion(rows)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, rows.Err()
}

func (r *Repository) Answer(ctx context.Context, tenantID, id, answer string) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE product_questions SET answer = $3, answered_at = NOW() WHERE tenant_id = $1 AND id = $2
	`, tenantID, id, answer)
	return err
}
