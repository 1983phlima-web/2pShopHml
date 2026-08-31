package domain

// SellerSummary powers the KPI cards + charts shown at the top of the
// seller's home screen: revenue, orders, units sold, a daily sales
// series, and the best-selling products — all computed from real orders,
// not mocked data.
type SellerSummary struct {
	TotalRevenue   int64          `json:"total_revenue"`
	TotalOrders    int            `json:"total_orders"`
	TotalUnitsSold int            `json:"total_units_sold"`
	PendingOrders  int            `json:"pending_orders"`
	SalesByDay     []DailySales   `json:"sales_by_day"`
	TopProducts    []ProductSales `json:"top_products"`
}

type DailySales struct {
	Date    string `json:"date"`
	Revenue int64  `json:"revenue"`
	Orders  int    `json:"orders"`
}

type ProductSales struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Units     int    `json:"units"`
	Revenue   int64  `json:"revenue"`
}

// AdminSummary powers the System/Global Admin dashboards: platform-wide
// counts by role/state/status plus total revenue.
type AdminSummary struct {
	UsersByRole     map[string]int `json:"users_by_role"`
	ProductsByState map[string]int `json:"products_by_state"`
	OrdersByStatus  map[string]int `json:"orders_by_status"`
	TotalRevenue    int64          `json:"total_revenue"`
	TotalReviews    int            `json:"total_reviews"`
}

// HealthReport powers the Global Admin health panel: infra-level signals
// beyond the basic /health liveness check.
type HealthReport struct {
	DatabaseOK        bool   `json:"database_ok"`
	MigrationsApplied int    `json:"migrations_applied"`
	TotalUsers        int    `json:"total_users"`
	TotalProducts     int    `json:"total_products"`
	TotalOrders       int    `json:"total_orders"`
	UptimeSeconds     int64  `json:"uptime_seconds"`
	Version           string `json:"version"`
}
