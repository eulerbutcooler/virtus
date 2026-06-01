package middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/internal/handler/response"
	"github.com/eulerbutcooler/virtus/internal/service"
	"github.com/google/uuid"
)

type ctxKey string

const (
	userIDKey   ctxKey = "user_id"
	userRoleKey ctxKey = "user_role"
)

// Wires the auth service into route protection.
type Authenticator struct {
	auth *service.AuthService
}

func NewAuthenticator(auth *service.AuthService) *Authenticator {
	return &Authenticator{auth: auth}
}

// Rejects requests without a valid Bearer access token.
func (a *Authenticator) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r)
		if !ok {
			response.Fail(w, http.StatusUnauthorized, "missing or malformed authorization header")
			return
		}
		claims, err := a.auth.ValidateAccessToken(token)
		if err != nil {
			response.Fail(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		ctx = context.WithValue(ctx, userRoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Allows the request through only if the authenticated user's role is in the allowed set.
func (a *Authenticator) RequireRole(roles ...domain.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := UserRoleFrom(r.Context())
			if slices.Contains(roles, role) {
				next.ServeHTTP(w, r)
				return
			}
			response.Fail(w, http.StatusForbidden, "forbidden")
		})
	}
}

func bearerToken(r *http.Request) (string, bool) {
	h := r.Header.Get("Authorization")
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}

// Returns the authenticated user ID set by RequireAuth.
func UserIDFrom(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(userIDKey).(uuid.UUID)
	return id
}

// Returns the authenticated user's role.
func UserRoleFrom(ctx context.Context) domain.UserRole {
	role, _ := ctx.Value(userRoleKey).(domain.UserRole)
	return role
}
