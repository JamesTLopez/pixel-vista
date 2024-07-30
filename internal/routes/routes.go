package routes

import (
	"embed"
	"fmt"
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

	router.Use(chiMiddle.Recoverer)
	router.Use(middleware.WithUser)
	router.Use(session.SessionManager.LoadAndSave)
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
		r.Get("/auth/callback", internal.GenerateHandler(handler.HandlerAuthCallback))
		r.Get("/login/provider/google", internal.GenerateHandler(handler.HandleLoginGoogleIndex))
		// Protected views
		r.Group(func(auth chi.Router) {
			auth.Use(middleware.WithAuth)
			auth.Get("/settings", internal.GenerateHandler(handler.HandleSettingsIndex))

		})
	})

	// endpoints
	router.Group(func(r chi.Router) {
		r.Post("/login", internal.GenerateHandler(handler.LoginCreate))
		r.Post("/register", internal.GenerateHandler(handler.RegisterCreate))
		r.Post("/logout", internal.GenerateHandler(handler.Logout))
	})

	return router
}

func getHandler(w http.ResponseWriter, r *http.Request) {

	msg := session.SessionManager.Keys(r.Context())

	fmt.Println("yeet-", msg)
	// io.WriteString(w, msg)
}
