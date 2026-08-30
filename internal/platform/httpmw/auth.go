package httpmw

import (
	"context"
	"net/http"
	"strings"

	"github.com/2pshop/2pshop/internal/identity/ports"
	"github.com/2pshop/2pshop/pkg/errors"
)

type claimsCtxKey struct{}

// RequireAuth validates the Bearer token and stores the resulting claims
// in the request context for downstream handlers and RequireRole.
func RequireAuth(tokens ports.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			token := strings.TrimPrefix(header, "Bearer ")
			if token == "" || token == header {
				writeError(w, errors.New(errors.ErrUnauthorized).WithDetail("reason", "missing bearer token"))
				return
			}
			claims, err := tokens.Validate(r.Context(), token)
			if err != nil {
				writeError(w, err)
				return
			}
			ctx := context.WithValue(r.Context(), claimsCtxKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole restricts access to the given set of roles. Must run after
// RequireAuth in the middleware chain.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFromContext(r.Context())
			if claims == nil || !allowed[claims.Role] {
				writeError(w, errors.New(errors.ErrForbidden).WithDetail("reason", "insufficient role"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ClaimsFromContext returns the authenticated user's token claims, or nil.
func ClaimsFromContext(ctx context.Context) *ports.TokenClaims {
	claims, _ := ctx.Value(claimsCtxKey{}).(*ports.TokenClaims)
	return claims
}
