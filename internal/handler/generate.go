package handler

import (
	"log/slog"
	"net/http"
	"pixelvista/db"
	"pixelvista/internal"
	"pixelvista/types"
	"pixelvista/view/pages/generate"

	"github.com/go-chi/chi/v5"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	user := internal.GetAuthenticatedUser(r)

	images, err := db.GetImagesByUserID(user.ID)

	if err != nil {
		return err
	}

	return renderComponent(w, r, generate.GeneratePage(generate.ViewData{
		Images: images,
	}))
}

func POSTGenerateImage(w http.ResponseWriter, r *http.Request) error {
	return renderComponent(w, r, generate.GalleryImage(types.Image{Status: types.ImageStatusPending}))
}

func GETGenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	slog.Info("checking image status", "id", id)
	image := types.Image{
		Status: types.ImageStatusCompleted,
	}

	return renderComponent(w, r, generate.GalleryImage(image))
}
