package routes

import (
	"embed"
	"net/http"
	"pixelvista/internal"
	"pixelvista/internal/handler"
	"pixelvista/internal/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddle "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func InitRoutes(FS embed.FS) http.Handler {
	router := chi.NewRouter()

	router.Use(chiMiddle.Recoverer)
	router.Use(middleware.WithUser)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"}, // TODO: for security, change such that it targets published origins
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           3000,
	}))

	// Allow styles in the public folder to be served
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))

	// Views
	router.Group(func(r chi.Router) {
		r.Get("/login", internal.GenerateHandler(handler.HandlerSigninIndex))
		r.Get("/dashboard", internal.GenerateHandler(handler.HandleHomeIndex))
		r.Get("/register", internal.GenerateHandler(handler.HandlerRegisterIndex))
	})

	// endpoints
	router.Group(func(r chi.Router) {
		r.Post("/login", internal.GenerateHandler(handler.LoginCreate))
	})

	return router
}
