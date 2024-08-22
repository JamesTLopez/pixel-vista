package handler

import (
	"log/slog"
	"net/http"
	"pixelvista/db"
	"pixelvista/internal"
	"pixelvista/types"
	"pixelvista/view/pages/generate"
	"strconv"

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
	user := internal.GetAuthenticatedUser(r)

	prompt := "eden in the garden"
	img := types.Image{
		Prompt: prompt,
		UserId: user.ID,
		Status: types.ImageStatusPending,
	}

	err := db.CreateImage(&img)

	if err != nil {
		return nil
	}

	return renderComponent(w, r, generate.GalleryImage(img))
}

func GETGenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		return err
	}
	slog.Info("checking image status", "id", id)

	image, err := db.GetImageById(id)

	if err != nil {
		return err
	}

	return renderComponent(w, r, generate.GalleryImage(image))
}
