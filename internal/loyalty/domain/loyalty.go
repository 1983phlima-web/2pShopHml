package domain

// Badge is one loyalty tier the user has earned in at least one period.
type Badge struct {
	Key      string `json:"key"` // "500" | "1000" | "vip"
	Label    string `json:"label"`
	EarnedAt string `json:"earned_at"`
}

// Profile is the customer's full loyalty snapshot: lifetime XP (1 XP per
// R$1 spent), accumulated Coins (from period badge awards), badges ever
// earned, and progress toward the next threshold in the current periods.
type Profile struct {
	XP                int64   `json:"xp"`
	Coins             int64   `json:"coins"`
	Badges            []Badge `json:"badges"`
	Period15DaySpend  int64   `json:"period_15day_spend"`
	Period15DayTarget int64   `json:"period_15day_target"`
	PeriodMonthSpend  int64   `json:"period_month_spend"`
	PeriodMonthTarget int64   `json:"period_month_target"`
}

// Thresholds, in cents.
const (
	Threshold500  int64 = 50_000
	Threshold1000 int64 = 100_000
	ThresholdVIP  int64 = 300_000
)

var BadgeLabels = map[string]string{
	"500":  "Bronze (R$500 em 15 dias)",
	"1000": "Prata (R$1.000 em 15 dias)",
	"vip":  "VIP (R$3.000 no mês)",
}
