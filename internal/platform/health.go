package platform

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type HealthChecker interface {
	Check(ctx context.Context) error
}

type HealthHandler struct {
	checks map[string]HealthChecker
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{checks: make(map[string]HealthChecker)}
}

func (h *HealthHandler) Register(name string, checker HealthChecker) {
	h.checks[name] = checker
}

func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	status := make(map[string]any)
	overall := "healthy"

	for name, checker := range h.checks {
		if err := checker.Check(ctx); err != nil {
			status[name] = map[string]any{"status": "unhealthy", "error": err.Error()}
			overall = "unhealthy"
		} else {
			status[name] = map[string]any{"status": "healthy"}
		}
	}

	status["status"] = overall
	status["timestamp"] = time.Now().UTC().Format(time.RFC3339)

	code := http.StatusOK
	if overall != "healthy" {
		code = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(status)
}

type DBHealthChecker struct {
	db *DB
}

func NewDBHealthChecker(db *DB) *DBHealthChecker {
	return &DBHealthChecker{db: db}
}

func (c *DBHealthChecker) Check(ctx context.Context) error {
	return c.db.Pool.Ping(ctx)
}
