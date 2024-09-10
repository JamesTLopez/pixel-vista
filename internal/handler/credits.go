package handler

import (
	"net/http"
	"pixelvista/view/pages/credits"
)

func HandleCreditsIndex(w http.ResponseWriter, r *http.Request) error {

	return renderComponent(w, r, credits.CreditsIndex())
}
