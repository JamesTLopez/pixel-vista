package handler

import (
	"net/http"
	"pixelvista/view/pages/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	return home.Index().Render(r.Context(), w)
}
