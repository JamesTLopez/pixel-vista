package handler

import (
	"net/http"
	"pixelvista/view/pages/generate"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	return renderComponent(w, r, generate.GeneratePage())
}
