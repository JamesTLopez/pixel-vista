package db

import (
	"context"
	"pixelvista/types"
)

func CreateAccount(account types.Account) error {
	_, err := Bun.NewInsert().
		Model(&account).
		Exec(context.Background())
	return err

}
