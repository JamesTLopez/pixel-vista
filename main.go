package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"pixelvista/internal/routes"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

type Application struct {
	Config Config
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	cfg := Config{
		Port: os.Getenv("PORT"),
	}

	app := &Application{
		Config: cfg,
	}

	err = app.Serve()

	if err != nil {
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
