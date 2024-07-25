package handler

import (
	"net/http"
	"pixelvista/helpers"
	"pixelvista/internal/sb"
	"pixelvista/view/pages/auth"

	"github.com/nedpals/supabase-go"
)

func HandlerSigninIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.Signin().Render(r.Context(), w)
}

func HandlerRegisterIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.Register().Render(r.Context(), w)
}

func LoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	if !helpers.IsValidEmail(credentials.Email) {

		return renderComponent(w, r, auth.LoginForm(credentials, auth.LoginErrors{
			Email: "Invalid email, please try again",
		}))
	}

	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)

	if err != nil {
		return renderComponent(w, r, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCred: "Invalid credentials, please try again",
		}))
	}

	cookie := &http.Cookie{
		Value:    resp.AccessToken,
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)
	hxRedirect(w, r, "/dashboard")
	return nil
}

// func RegisterCreate(w http.ResponseWriter, r *http.Request) error {
// 	params := auth.RegisterParams{
// 		Email:           r.FormValue("email"),
// 		Password:        r.FormValue("password"),
// 		ConfirmPassword: r.FormValue("confirmPassword"),
// 	}

// }
