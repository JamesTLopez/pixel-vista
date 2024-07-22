package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"pixelvista/router"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

type Application struct {
	Config Config
}

func (app *Application) Serve() error {
	port := os.Getenv("PORT")
	slog.Info("Started application...", "port", port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router.Routes(),
	}

	return server.ListenAndServe()
}

func main() {
	if err := initPixelVista(); err != nil {
		log.Fatal(err)
	}

	cfg := Config{
		Port: os.Getenv("PORT"),
	}

	app := &Application{
		Config: cfg,
	}

	err := app.Serve()

	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func initPixelVista() error {
	err := godotenv.Load()

	if err != nil {
		return err
	}

	return nil
}
