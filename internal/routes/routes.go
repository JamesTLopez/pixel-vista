package routes

import (
	"embed"
	"net/http"
	"os"
	"pixelvista/internal"
	"pixelvista/internal/handler"
	"pixelvista/internal/middleware"
	"pixelvista/internal/session"
	"strings"

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

	allowed_origins := os.Getenv("ALLOWED_ORIGINS")

	origins := strings.Split(allowed_origins, ",")

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
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
		r.Post("/replicate/callback/{userID}/{batchID}", internal.GenerateHandler(handler.ReplicateCallback))

		r.Group(func(r chi.Router) {
			r.Use(middleware.WithAccountSetup)
			r.Get("/", internal.GenerateHandler(handler.HandleHomeIndex)) // Maybe ?
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.WithAccountSetup, middleware.WithAuth)
			r.Get("/settings", internal.GenerateHandler(handler.HandleSettingsIndex))
			r.Put("/settings/account/profile", internal.GenerateHandler(handler.HandleSettingsProfileUpdate))
			r.Get("/generate", internal.GenerateHandler(handler.HandleGenerateIndex))
			r.Post("/generate", internal.GenerateHandler(handler.POSTGenerateImage))
			r.Get("/generate/image/status/{id}", internal.GenerateHandler(handler.GETGenerateImageStatus))
			r.Get("/buy-credits", internal.GenerateHandler(handler.HandleCreditsIndex))
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.WithAuth)
			r.Post("/account/setup", internal.GenerateHandler(handler.SetupAccountCreate))
			r.Get("/account/setup", internal.GenerateHandler(handler.HandlerAccountIndex))
		})
	})
	return router
}
