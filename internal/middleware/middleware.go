package middleware

import (
	"context"
	"net/http"
	"pixelvista/internal"
	"pixelvista/types"
	"strings"
)

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		user := internal.GetAuthenticatedUser(r)
		// user := types.AuthenticatedUser{
		// 	Email:    "test1@gmail.com",
		// 	LoggedIn: true,
		// }
		ctx := context.WithValue(r.Context(), types.Userkey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
