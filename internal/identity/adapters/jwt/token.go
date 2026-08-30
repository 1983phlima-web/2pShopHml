// Package jwt provides a JWT-based implementation of ports.TokenService.
package jwt

import (
	"context"
	"time"

	"github.com/2pshop/2pshop/internal/identity/domain"
	"github.com/2pshop/2pshop/internal/identity/ports"
	"github.com/2pshop/2pshop/pkg/errors"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenService(secret string, ttl time.Duration) *TokenService {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &TokenService{secret: []byte(secret), ttl: ttl}
}

type claims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwtlib.RegisteredClaims
}

func (t *TokenService) Generate(ctx context.Context, user *domain.User) (string, error) {
	now := time.Now().UTC()
	c := claims{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Email:    user.Email,
		Role:     string(user.Role),
		RegisteredClaims: jwtlib.RegisteredClaims{
			IssuedAt:  jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(now.Add(t.ttl)),
			Subject:   user.ID,
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, c)
	return token.SignedString(t.secret)
}

func (t *TokenService) Validate(ctx context.Context, tokenString string) (*ports.TokenClaims, error) {
	var c claims
	token, err := jwtlib.ParseWithClaims(tokenString, &c, func(tok *jwtlib.Token) (any, error) {
		return t.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New(errors.ErrUnauthorized).WithDetail("reason", "invalid token")
	}
	return &ports.TokenClaims{
		UserID:   c.UserID,
		TenantID: c.TenantID,
		Email:    c.Email,
		Role:     c.Role,
	}, nil
}
