package ports

import (
	"context"
	"time"
)

type Award struct {
	Badge    string
	Coins    int
	EarnedAt time.Time
}

type Repository interface {
	// LifetimeSpend sums all order totals ever placed by the user (cents).
	LifetimeSpend(ctx context.Context, tenantID, userID string) (int64, error)
	// PeriodSpend sums order totals placed since the given time (cents).
	PeriodSpend(ctx context.Context, tenantID, userID string, since time.Time) (int64, error)
	// GrantAward idempotently records a badge/coins award for a period —
	// a UNIQUE constraint on (tenant,user,period_type,period_key,badge)
	// guarantees each period only ever grants a given badge once.
	GrantAward(ctx context.Context, tenantID, userID, periodType, periodKey, badge string, coins int) error
	ListAwards(ctx context.Context, tenantID, userID string) ([]Award, error)
	TotalCoins(ctx context.Context, tenantID, userID string) (int64, error)
}
