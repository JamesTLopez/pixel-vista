package handler

import (
	"net/http"
	"pixelvista/view/pages/auth"

	"github.com/a-h/templ"
	"github.com/nedpals/supabase-go"
)

func renderComponent(w http.ResponseWriter, r *http.Request, component templ.Component) error {
	return component.Render(r.Context(), w)
}

func HandlerSigninIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.Signin().Render(r.Context(), w)
}

func LoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	return renderComponent(w, r, auth.LoginForm(credentials, auth.LoginErrors{
		InvalidCred: "Invalid credentials, please try again",
	}))
}
