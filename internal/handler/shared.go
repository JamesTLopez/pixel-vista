package handler

import (
	"net/http"

	"github.com/a-h/templ"
)

func renderComponent(w http.ResponseWriter, r *http.Request, component templ.Component) error {
	return component.Render(r.Context(), w)
}
func hxRedirect(w http.ResponseWriter, r *http.Request, to string) error {
	if len(r.Header.Get("HX-Request")) > 0 {
		w.Header().Set("HX-Redirect", to)
		w.WriteHeader(http.StatusSeeOther)
		return nil
	}
	http.Redirect(w, r, to, http.StatusSeeOther)
	return nil
}
