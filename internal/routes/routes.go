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

	allowed_origins := os.Getenv("ALLOWED_ORIGINS")

	origins := strings.Split(allowed_origins, ",")

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           3000,
	}))

	// Custom middleware
	router.Use(middleware.WithUser)

	// Allow styles in the public folder to be served
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))

	router.Get("/", internal.GenerateHandler(handler.HandleHomeIndex)) // Maybe ?
	router.Get("/login", internal.GenerateHandler(handler.HandlerSigninIndex))
	router.Get("/login/provider/google", internal.GenerateHandler(handler.HandleLoginGoogleIndex))
	router.Get("/auth/callback", internal.GenerateHandler(handler.HandlerAuthCallback))
	router.Post("/logout", internal.GenerateHandler(handler.Logout))
	router.Post("/login", internal.GenerateHandler(handler.LoginCreate))
	router.Post("/replicate/callback/{userID}/{batchID}", internal.GenerateHandler(handler.ReplicateCallback))

	router.Group(func(r chi.Router) {
		r.Use(middleware.WithAuth, middleware.WithAccountSetup)
		r.Get("/settings", internal.GenerateHandler(handler.HandleSettingsIndex))
		r.Put("/settings/account/profile", internal.GenerateHandler(handler.HandleSettingsProfileUpdate))
		r.Get("/generate", internal.GenerateHandler(handler.HandleGenerateIndex))
		r.Post("/generate", internal.GenerateHandler(handler.POSTGenerateImage))
		r.Get("/generate/image/status/{id}", internal.GenerateHandler(handler.GETGenerateImageStatus))
		r.Get("/buy-credits", internal.GenerateHandler(handler.HandleCreditsIndex))
		r.Get("/checkout/create/{productID}", internal.GenerateHandler(handler.StripeCheckout))
		r.Get("/checkout/success/{sessionID}", internal.GenerateHandler(handler.StripeCheckoutSuccess))
		r.Get("/checkout/cancel", internal.GenerateHandler(handler.StripeCheckoutCancel))

	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.WithAuth)
		r.Post("/account/setup", internal.GenerateHandler(handler.SetupAccountCreate))
		r.Get("/account/setup", internal.GenerateHandler(handler.HandlerAccountIndex))
	})
	return router
}
