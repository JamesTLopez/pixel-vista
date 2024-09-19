package db

import (
	"context"
	"pixelvista/types"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func CreateImage(bun bun.Tx, image *types.Image) error {
	_, err := bun.NewInsert().Model(image).Exec(context.Background())

	return err
}

func UpdateImage(bun bun.Tx, image *types.Image) error {

	_, err := bun.NewUpdate().Model(image).WherePK().Exec(context.Background())

	return err
}

func GetImageById(id int) (types.Image, error) {
	var image types.Image

	err := Bun.NewSelect().Model(&image).Where("id = ?", id).Scan(context.Background())

	return image, err
}

func GetImagesByUserID(userID uuid.UUID) ([]types.Image, error) {
	var images []types.Image

	err := Bun.NewSelect().Model(&images).Where("deleted = ?", false).Where("user_id = ?", userID).Scan(context.Background())

	return images, err
}

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

func UpdateProfile(account *types.Account) error {
	_, err := Bun.NewUpdate().Model(account).WherePK().Exec(context.Background())
	return err
}

func GetImagesByBatchID(batchID uuid.UUID) ([]types.Image, error) {
	var images []types.Image

	err := Bun.NewSelect().Model(&images).Where("batch_id = ?", batchID).Scan(context.Background())

	return images, err

}

func GetCreditPrices() ([]types.CreditPrice, error) {
	var creditPrice []types.CreditPrice

	err := Bun.NewSelect().Model(&creditPrice).Order("created_at ASC").Scan(context.Background())

	return creditPrice, err
}

func GetCreditPriceByID(productID string) (types.CreditPrice, error) {
	var price types.CreditPrice

	err := Bun.NewSelect().Model(&price).Where("product_id = ?", productID).Scan(context.Background())

	return price, err
}
