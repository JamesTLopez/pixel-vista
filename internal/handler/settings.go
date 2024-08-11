package handler

import (
	"net/http"
	"pixelvista/db"
	"pixelvista/internal"
	"pixelvista/pkg/validation"
	settings "pixelvista/view/pages/settings"
)

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := internal.GetAuthenticatedUser(r)
	return renderComponent(w, r, settings.SettingsIndex(user))
}

func HandleSettingsProfileUpdate(w http.ResponseWriter, r *http.Request) error {

	params := settings.UpdateProfileParams{
		Username: r.FormValue("username"),
	}
	errors := settings.UpdateProfileErrors{}
	ok := validation.New(&params, validation.Fields{
		"Username": validation.Rules(validation.Min(3), validation.Max(20)),
	}).Validate(&errors)

	user := internal.GetAuthenticatedUser(r)
	if !ok {
		return renderComponent(w, r, settings.SettingsProfileForm(params, errors, user))
	}

	user.Account.Username = params.Username

	if err := db.UpdateProfile(&user.Account); err != nil {
		return err
	}

	params.Success = true

	return renderComponent(w, r, settings.SettingsProfileForm(params, settings.UpdateProfileErrors{}, user))
}
