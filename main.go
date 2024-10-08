package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"pixelvista/db"
	replicate "pixelvista/internal/repl"
	"pixelvista/internal/routes"
	superb "pixelvista/internal/sb"
	"pixelvista/internal/session"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	cfg := Config{
		Port: os.Getenv("PORT"),
	}

	app := &Application{
		Config: cfg,
	}

	session.InitSession()
	superb.SbInit()

	err := replicate.ReplicateInit()

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Init(); err != nil {
		log.Fatal(err)

	}
	if err := app.Serve(); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

//go:embed public
var FS embed.FS

func (app *Application) Serve() error {
	port := os.Getenv("PORT")
	slog.Info("Started application...", "port", port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: routes.InitRoutes(FS),
	}

	return server.ListenAndServe()
}

type Config struct {
	Port string
}

type Application struct {
	Config Config
}
