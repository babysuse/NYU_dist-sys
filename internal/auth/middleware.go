package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/jwt"
	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/users"
)

type contextKey string

var userContext = contextKey("user")

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// allow unauthenticated user
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			// validate jwt token
			tokenStr := header
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				return
			}

			// get the user from the database
			user := users.User{Username: username}
			id, err := users.GetUserIdByUsername(username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// put it in context
			user.ID = strconv.Itoa(id)
			ctx := context.WithValue(r.Context(), userContext, &user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(userContext).(*users.User)
	return raw
}
