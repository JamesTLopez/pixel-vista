package middleware

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"pixelvista/db"
	"pixelvista/internal"
	"pixelvista/internal/sb"
	"pixelvista/internal/session"
	"pixelvista/types"
	"strings"

	"github.com/google/uuid"
)

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		_, err := r.Cookie("pixel_vista")

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		access_token := session.SessionManager.GetString(r.Context(), "accessToken")

		resp, err := sb.Client.Auth.User(r.Context(), access_token)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user := types.AuthenticatedUser{
			ID:       uuid.MustParse(resp.ID),
			Email:    resp.Email,
			LoggedIn: true,
		}

		ctx := context.WithValue(r.Context(), types.Userkey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
func WithAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		user := internal.GetAuthenticatedUser(r)
		if !user.LoggedIn {
			path := r.URL.Path
			http.Redirect(w, r, "/login?to="+path, http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func WithAccountSetup(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := internal.GetAuthenticatedUser(r)

		if !user.LoggedIn {
			next.ServeHTTP(w, r)
			return
		}

		account, err := db.GetAccountGyUserID(user.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Redirect(w, r, "/account/setup", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		user.Account = account
		ctx := context.WithValue(r.Context(), types.Userkey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
