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
	ID       int `bun:"id,pk,autoincrement"`
	UserId   uuid.UUID
	Status   ImageStatus
	CreateAt time.Time `bun:"default:'now()'"`
}
