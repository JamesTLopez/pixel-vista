package db

import (
	"errors"
	"os"

	supa "github.com/nedpals/supabase-go"
)

var (
	Client   *supa.Client
	sbUrl    string
	sbSecret string
)

func SbInit() error {

	if sbUrl := os.Getenv("SUPABASE_URL"); sbUrl == "" {
		return errors.New("supabase host is required")
	}

	if sbSecret := os.Getenv("SUPABASE_SECRET"); sbSecret == "" {
		return errors.New("supabase secret not provided")
	}

	Client = supa.CreateClient(sbUrl, sbSecret)

	return nil
}
