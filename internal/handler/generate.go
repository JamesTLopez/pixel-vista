package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"pixelvista/db"
	"pixelvista/internal"
	repl "pixelvista/internal/repl"
	"pixelvista/pkg/validation"
	"pixelvista/types"
	"pixelvista/view/pages/generate"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/replicate/replicate-go"
	"github.com/uptrace/bun"
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
	amount, _ := strconv.Atoi(r.FormValue("amount"))

	params := generate.FormParams{
		Prompt: r.FormValue("prompt"),
		Amount: amount,
	}

	var errors generate.FormErrors
	ok := validation.New(params, validation.Fields{
		"Prompt": validation.Rules(validation.Min(10), validation.Max(200)),
	}).Validate(&errors)

	if amount <= 0 || amount > 8 {
		errors.Amount = "Please enter a valid amount"
		return renderComponent(w, r, generate.GenerateForm(params, errors))
	}
	if !ok {
		return renderComponent(w, r, generate.GenerateForm(params, errors))
	}

	batchID := uuid.New()
	genImageParams := GenerateImageParams{
		Prompt:  params.Prompt,
		Amount:  params.Amount,
		UserID:  user.ID,
		BatchID: batchID,
	}

	if err := generateImage(r.Context(), genImageParams); err != nil {
		return err
	}

	return db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		for x := 0; x < amount; x++ {
			img := types.Image{
				Prompt:  params.Prompt,
				UserId:  user.ID,
				Status:  types.ImageStatusPending,
				BatchID: batchID,
			}

			if err := db.CreateImage(&img); err != nil {
				return err
			}
		}
		// TODO: use hx-target instead of redirect on success
		return hxRedirect(w, r, "/generate")
	})

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

type GenerateImageParams struct {
	Prompt  string
	Amount  int
	BatchID uuid.UUID
	UserID  uuid.UUID
}

func generateImage(ctx context.Context, params GenerateImageParams) error {
	input := replicate.PredictionInput{
		"prompt":      params.Prompt,
		"num_outputs": params.Amount,
	}

	webhookURL := os.Getenv("WEBHOOK")

	webhook := replicate.Webhook{
		URL:    fmt.Sprintf("%s/replicate/callback/%s/%s", webhookURL, params.UserID, params.BatchID),
		Events: []replicate.WebhookEventType{"start", "completed"},
	}

	_, err := repl.ReplicateClient.CreatePredictionWithModel(ctx, "black-forest-labs", "flux-schnell", input, &webhook, false)

	return err
}
