package auth

import (
	"context"
	"net/http"
)

type contextKey string

const UserContextKey contextKey = "user"

type UserContext struct {
	UserID string
	Email  string
}

type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := ValidateJWT(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userCtx := &UserContext{
			UserID: claims.UserID,
			Email:  claims.Email,
		}
		ctx := context.WithValue(r.Context(), UserContextKey, userCtx)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) (*UserContext, bool) {
	user, ok := ctx.Value(UserContextKey).(*UserContext)
	return user, ok
}