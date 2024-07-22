package routes

import (
	"embed"
	"net/http"
	"pixelvista/handler"
	"pixelvista/helpers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func InitRoutes(FS embed.FS) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"}, // TODO: for security, change such that it targets published origins
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           3000,
	}))

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))

	// Views
	router.Group(func(r chi.Router) {
		r.Get("/", helpers.GenerateHandler(handler.HandleHomeIndex))
	})

	return router
}
