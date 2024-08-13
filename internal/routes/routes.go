package routes

import (
	"embed"
	"net/http"
	"pixelvista/internal"
	"pixelvista/internal/handler"
	"pixelvista/internal/middleware"
	"pixelvista/internal/session"

	"github.com/go-chi/chi/v5"
	chiMiddle "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func InitRoutes(FS embed.FS) http.Handler {
	router := chi.NewRouter()

	// Define 3rd party library first
	router.Use(chiMiddle.Recoverer)
	router.Use(session.SessionManager.LoadAndSave)

	// Custom middleware
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

	router.Group(func(r chi.Router) {
		r.Get("/login", internal.GenerateHandler(handler.HandlerSigninIndex))
		r.Get("/login/provider/google", internal.GenerateHandler(handler.HandleLoginGoogleIndex))
		r.Get("/auth/callback", internal.GenerateHandler(handler.HandlerAuthCallback))
		r.Post("/logout", internal.GenerateHandler(handler.Logout))
		r.Post("/login", internal.GenerateHandler(handler.LoginCreate))

		r.Group(func(r chi.Router) {
			r.Use(middleware.WithAccountSetup)
			r.Get("/", internal.GenerateHandler(handler.HandleHomeIndex)) // Maybe ?
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.WithAccountSetup, middleware.WithAuth)
			r.Get("/settings", internal.GenerateHandler(handler.HandleSettingsIndex))
			r.Put("/settings/account/profile", internal.GenerateHandler(handler.HandleSettingsProfileUpdate))
			r.Get("/generate", internal.GenerateHandler(handler.HandleGenerateIndex))
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.WithAuth)
			r.Post("/account/setup", internal.GenerateHandler(handler.SetupAccountCreate))
			r.Get("/account/setup", internal.GenerateHandler(handler.HandlerAccountIndex))
		})
	})
	return router
}
