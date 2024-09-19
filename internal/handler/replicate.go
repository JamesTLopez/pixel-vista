package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"pixelvista/db"
	"pixelvista/types"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/go-chi/chi/v5"
)

const (
	succeeded  = "succeeded"
	processing = "processing"
)

func ReplicateCallback(w http.ResponseWriter, r *http.Request) error {
	var resp ReplicateResponse

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return err
	}

	if resp.Status == processing {
		return nil
	}

	if resp.Status != succeeded {
		return fmt.Errorf("replicate callback responded with a non ok status %s", resp.Status)
	}

	batchID, err := uuid.Parse(chi.URLParam(r, "batchID"))

	if err != nil {
		return fmt.Errorf("replicate callback has invalid batchID: %s", err)
	}

	images, err := db.GetImagesByBatchID(batchID)
	if err != nil {
		return fmt.Errorf("replicate callback has failed to find images with batchID: %s error=%s", batchID, err)
	}

	if len(images) != len(resp.Output) {
		return fmt.Errorf("replicate callback un-equal image compared to replicate response")
	}

	db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for i, imageURl := range resp.Output {
			images[i].Status = types.ImageStatusCompleted
			images[i].ImageUrl = imageURl
			images[i].Prompt = resp.Input.Prompt
			if err := db.UpdateImage(tx, &images[i]); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

type Input struct {
	NumOutputs int    `json:"num_ouputs"`
	Prompt     string `json:"prompt"`
}

type Metrics struct {
	ImageCount  int     `json:"image_count"`
	PredictTime float64 `json:"predict_time"`
}

type URLs struct {
	Cancel string `json:"cancel"`
	Get    string `json:"get"`
}

type ReplicateResponse struct {
	CompletedAt         time.Time `json:"completed_at"`
	CreatedAt           time.Time `json:"created_at"`
	DataRemoved         bool      `json:"data_removed"`
	Error               *string   `json:"error"`
	ID                  string    `json:"id"`
	Input               Input     `json:"input"`
	Logs                string    `json:"logs"`
	Metrics             Metrics   `json:"metrics"`
	Model               string    `json:"model"`
	Output              []string  `json:"output"`
	StartedAt           time.Time `json:"started_at"`
	Status              string    `json:"status"`
	URLs                URLs      `json:"urls"`
	Version             string    `json:"version"`
	Webhook             string    `json:"webhook"`
	WebhookEventsFilter []string  `json:"webhook_events_filter"`
}
