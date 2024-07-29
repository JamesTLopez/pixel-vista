package handler

import (
	"net/http"
	"pixelvista/internal"
	settings "pixelvista/view/pages/settings"
)

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := internal.GetAuthenticatedUser(r)
	return renderComponent(w, r, settings.SettingsIndex(user))
}
