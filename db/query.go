package db

import (
	"context"
	"pixelvista/types"

	"github.com/google/uuid"
)

func CreateAccount(account types.Account) error {
	_, err := Bun.NewInsert().
		Model(&account).
		Exec(context.Background())
	return err

}

func GetAccountGyUserID(id uuid.UUID) (types.Account, error) {
	var account types.Account

	err := Bun.NewSelect().Model(&account).Where("user_id = ?", id).Scan(context.Background())
	return account, err
}
