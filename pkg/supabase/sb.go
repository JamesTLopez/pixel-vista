package supabase

import (
	"errors"
	"os"
)

var (
	// Client   *supa.Client
	sbHost   string
	sbSecret string
)

func SbInit() error {
	sbUrl := os.Getenv("SUPABASE_URL")

	if sbUrl == "" {
		return errors.New("Supabase host is required")
	}

	if sbSecret := os.Getenv("SUPABASE_SECRET"); sbSecret == "" {
		return errors.New("Supabase secret not provided")
	}

	// Client = supa.CreateClient(sbUrl, sbSecret)

	return nil
}
