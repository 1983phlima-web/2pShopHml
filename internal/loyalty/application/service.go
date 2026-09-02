package application

import (
	"context"
	"fmt"
	"time"

	"github.com/2pshop/2pshop/internal/loyalty/domain"
	"github.com/2pshop/2pshop/internal/loyalty/ports"
	"github.com/2pshop/2pshop/pkg/errors"
)

type Service struct {
	repo ports.Repository
}

func NewService(repo ports.Repository) *Service {
	return &Service{repo: repo}
}

// currentPeriods computes the current fixed 15-day window (1st-15th or
// 16th-end of month) and the current calendar month window, each with a
// stable period key used as the idempotency key for badge grants.
func currentPeriods(now time.Time) (period15Start time.Time, period15Key string, monthStart time.Time, monthKey string) {
	year, month, day := now.Date()
	if day <= 15 {
		period15Start = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		period15Key = fmt.Sprintf("%04d-%02d-P1", year, month)
	} else {
		period15Start = time.Date(year, month, 16, 0, 0, 0, 0, time.UTC)
		period15Key = fmt.Sprintf("%04d-%02d-P2", year, month)
	}
	monthStart = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	monthKey = fmt.Sprintf("%04d-%02d", year, month)
	return
}

// GetProfile computes the customer's live loyalty snapshot: it reconciles
// (idempotently) any badge thresholds crossed in the current periods, then
// returns lifetime XP, accumulated coins, and every badge ever earned.
func (s *Service) GetProfile(ctx context.Context, tenantID, userID string) (*domain.Profile, error) {
	now := time.Now().UTC()
	period15Start, period15Key, monthStart, monthKey := currentPeriods(now)

	period15Spend, err := s.repo.PeriodSpend(ctx, tenantID, userID, period15Start)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to compute 15-day spend", err)
	}
	monthSpend, err := s.repo.PeriodSpend(ctx, tenantID, userID, monthStart)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to compute month spend", err)
	}

	// Reconcile: grant whichever 15-day tier applies (idempotent — a
	// period already awarded is simply skipped by the UNIQUE constraint).
	if period15Spend >= domain.Threshold1000 {
		_ = s.repo.GrantAward(ctx, tenantID, userID, "15day", period15Key, "1000", 1000)
	} else if period15Spend >= domain.Threshold500 {
		_ = s.repo.GrantAward(ctx, tenantID, userID, "15day", period15Key, "500", 500)
	}
	if monthSpend >= domain.ThresholdVIP {
		_ = s.repo.GrantAward(ctx, tenantID, userID, "month", monthKey, "vip", 3000)
	}

	lifetimeSpend, err := s.repo.LifetimeSpend(ctx, tenantID, userID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to compute lifetime spend", err)
	}
	coins, err := s.repo.TotalCoins(ctx, tenantID, userID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to compute total coins", err)
	}
	awards, err := s.repo.ListAwards(ctx, tenantID, userID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to list awards", err)
	}

	seen := make(map[string]bool)
	var badges []domain.Badge
	for _, a := range awards {
		if seen[a.Badge] {
			continue
		}
		seen[a.Badge] = true
		badges = append(badges, domain.Badge{
			Key:      a.Badge,
			Label:    domain.BadgeLabels[a.Badge],
			EarnedAt: a.EarnedAt.Format(time.RFC3339),
		})
	}

	return &domain.Profile{
		XP:                lifetimeSpend / 100, // cents -> reais, 1:1 with XP
		Coins:             coins,
		Badges:            badges,
		Period15DaySpend:  period15Spend,
		Period15DayTarget: domain.Threshold1000,
		PeriodMonthSpend:  monthSpend,
		PeriodMonthTarget: domain.ThresholdVIP,
	}, nil
}
