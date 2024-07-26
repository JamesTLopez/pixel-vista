package handler

import (
	"fmt"
	"net/http"
	"pixelvista/helpers/validation"
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
	if !validation.IsValidEmail(credentials.Email) {

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

func RegisterCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.RegisterParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}
	fmt.Println(params)
	errors := auth.RegisterErrors{}

	if ok := validation.New(&params, validation.Fields{
		"Email":           validation.Rules(validation.Email),
		"Password":        validation.Rules(validation.Password),
		"ConfirmPassword": validation.Rules(validation.Equal(params.Password), validation.Message("Passwords do not match")),
	}).Validate(&errors); !ok {
		return renderComponent(w, r, auth.RegisterForm(params, errors))
	}

	resp, err := sb.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{Email: params.Email, Password: params.Password})

	if err != nil {
		return err
	}

	return renderComponent(w, r, auth.RegisterSuccess(resp.Email))
}
