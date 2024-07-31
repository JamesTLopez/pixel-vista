package handler

import (
	"fmt"
	"net/http"
	"pixelvista/db"
	"pixelvista/internal"
	"pixelvista/types"
	"pixelvista/view/pages/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {

	user := internal.GetAuthenticatedUser(r)
	account := types.Account{
		UserID:   user.ID,
		Username: "foo",
	}

	if err := db.CreateAccount(account); err != nil {
		return err
	}

	fmt.Printf("%+v\n", account)

	return home.Index().Render(r.Context(), w)
}
