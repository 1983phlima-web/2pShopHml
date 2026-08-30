// Package httpmw provides cross-cutting HTTP middleware shared across
// all API modules: tenant resolution and authentication/RBAC.
package httpmw

import (
	"context"
	"encoding/json"
	"net/http"

	tenancyApp "github.com/2pshop/2pshop/internal/tenancy/application"
	"github.com/2pshop/2pshop/pkg/errors"
)

type tenantCtxKey struct{}

// TenantResolver reads the X-Tenant-ID header (a tenant slug) and resolves
// it to the tenant's UUID, making it available via TenantFromContext.
func TenantResolver(tenancy *tenancyApp.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slug := r.Header.Get("X-Tenant-ID")
			if slug == "" {
				writeError(w, errors.New(errors.ErrInvalidInput).WithDetail("header", "X-Tenant-ID is required"))
				return
			}
			tenant, err := tenancy.GetTenantBySlug(r.Context(), slug)
			if err != nil {
				writeError(w, err)
				return
			}
			ctx := context.WithValue(r.Context(), tenantCtxKey{}, tenant.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TenantFromContext returns the resolved tenant UUID, or "" if absent.
func TenantFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(tenantCtxKey{}).(string); ok {
		return v
	}
	return ""
}

func writeError(w http.ResponseWriter, err error) {
	status := errors.HTTPStatus(err)
	var appErr *errors.AppError
	if ae, ok := err.(*errors.AppError); ok {
		appErr = ae
	} else {
		appErr = errors.New(errors.ErrInternal)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(appErr)
}
