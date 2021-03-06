package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/os3224/final-project-0b5a2e16-babysuse/internal/pkg/jwt"
	"github.com/os3224/final-project-0b5a2e16-babysuse/web/account/pb"
)

type contextKey string

var userContext = contextKey("user")

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := r.Header.Get("Authorization")

			fmt.Printf("token: %s\n", tokenStr)
			// allow unauthenticated user
			if tokenStr == "" {
				next.ServeHTTP(w, r)
				return
			}

			// validate jwt token
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				return
			}

			// get the user from the database
			user := pb.Account{Username: username}
			//id, err := account.GetUserIdByUsername(username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// put it in context
			//user.ID = strconv.Itoa(id)
			ctx := context.WithValue(r.Context(), userContext, &user)
			fmt.Printf("%s logged in\n", user.Username)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// find user from context
func ForContext(ctx context.Context) *pb.Account {
	raw, _ := ctx.Value(userContext).(*pb.Account)
	return raw
}
