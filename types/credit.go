package types

import (
	"time"
)

type CreditPrice struct {
	ID        int `bun:"id,pk,autoincrement"`
	ProductId string
	Name      string
	Price     string
	Credits   int
	CreatedAt time.Time `bun:"default:'now()'"`
}
