package handler

import (
	"net/http"
	"pixelvista/view/home"
)


func HandleHomeIndex(w http.ResponseWriter, r *http.Request) {
	home.Index().Render(r.Context(), w)
}