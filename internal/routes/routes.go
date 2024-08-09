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

	// Views
	router.Group(func(r chi.Router) {
		r.Get("/login", internal.GenerateHandler(handler.HandlerSigninIndex))
		r.Get("/register", internal.GenerateHandler(handler.HandlerRegisterIndex))
		r.Get("/auth/callback", internal.GenerateHandler(handler.HandlerAuthCallback))
		r.Get("/login/provider/google", internal.GenerateHandler(handler.HandleLoginGoogleIndex))

		// Protected views
		r.Group(func(acc chi.Router) {
			acc.Use(middleware.WithAccountSetup, middleware.WithAuth)
			acc.Get("/", internal.GenerateHandler(handler.HandleHomeIndex))
			acc.Get("/account/setup", internal.GenerateHandler(handler.HandlerAccountIndex))
			acc.Get("/settings", internal.GenerateHandler(handler.HandleSettingsIndex))
		})
	})

	// endpoints
	router.Group(func(r chi.Router) {
		r.Use(middleware.WithAccountSetup, middleware.WithAuth)
		r.Post("/login", internal.GenerateHandler(handler.LoginCreate))
		r.Post("/register", internal.GenerateHandler(handler.RegisterCreate))
		r.Post("/logout", internal.GenerateHandler(handler.Logout))
		r.Post("/account/setup", internal.GenerateHandler(handler.SetupAccountCreate))
		r.Put("/settings/account/profile", internal.GenerateHandler(handler.HandleSettingsProfileUpdate))
	})

	return router
}
