package ports

import (
	"context"

	catalogDomain "github.com/2pshop/2pshop/internal/catalog/domain"
)

type Repository interface {
	Add(ctx context.Context, tenantID, userID, productID string) error
	Remove(ctx context.Context, tenantID, userID, productID string) error
	ListProductIDs(ctx context.Context, tenantID, userID string) ([]string, error)
	// ListWithProducts returns the user's favorited products with full
	// catalog data (joined), so the frontend can render them with the
	// same product card used across the storefront.
	ListWithProducts(ctx context.Context, tenantID, userID string) ([]*catalogDomain.Product, error)
}
