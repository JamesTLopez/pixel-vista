package helpers

import (
	"log/slog"
	"net/http"
)

func GenerateHandler(handler func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			slog.Error("Internal server error", "err", err, "path", r.URL.Path)
		}
	}
}
