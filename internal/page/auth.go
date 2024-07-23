package page

import (
	"net/http"
	"pixelvista/view/pages/auth"
)

func HandlerSigninIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.Signin().Render(r.Context(), w)
}
