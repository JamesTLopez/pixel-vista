package types

import (
	"time"

	"github.com/google/uuid"
)

type ImageStatus int

const (
	ImageStatusPending ImageStatus = iota
	ImageStatusFailed
	ImageStatusCompleted
)

type Image struct {
	ID        int `bun:"id,pk,autoincrement"`
	UserId    uuid.UUID
	Status    ImageStatus
	Prompt    string
	ImageUrl  string
	deleted   bool      `bun:"default:'false'"`
	CreatedAt time.Time `bun:"default:'now()'"`
	DeletedAt time.Time
}
