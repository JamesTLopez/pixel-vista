package internal

import (
	"log/slog"
	"net/http"
	"pixelvista/types"
)

func GetAuthenticatedUser(r *http.Request) types.AuthenticatedUser {
	user, ok := r.Context().Value(types.Userkey).(types.AuthenticatedUser)

	if !ok {
		return types.AuthenticatedUser{}
	}

	return user
}

func GenerateHandler(handler func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			slog.Error("Internal server error", "err", err, "path", r.URL.Path)
		}
	}
}
