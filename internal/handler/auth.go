package handler

import (
	"net/http"
	"os"
	"pixelvista/db"
	"pixelvista/internal"
	"pixelvista/internal/sb"
	"pixelvista/internal/session"
	"pixelvista/pkg/validation"
	"pixelvista/types"
	"pixelvista/view/pages/auth"

	"github.com/nedpals/supabase-go"
)

func HandlerSigninIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.Signin().Render(r.Context(), w)
}

func HandlerAccountIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.AccountCreationPage().Render(r.Context(), w)
}

func SetupAccountCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.AccountSetupFormParams{
		Username: r.FormValue("username"),
	}

	var errors auth.AccountSetupFormError
	fields := validation.Fields{
		"Username": validation.Rules(validation.Min(2), validation.Max(20)),
	}

	ok := validation.New(&params, fields).Validate(&errors)

	if !ok {
		return renderComponent(w, r, auth.AccountSetupForm(params, errors))
	}

	user := internal.GetAuthenticatedUser(r)
	account := types.Account{
		UserID:   user.ID,
		Username: params.Username,
	}

	if err := db.CreateAccount(account); err != nil {
		return err
	}

	return hxRedirect(w, r, "/")
}

func HandleLoginGoogleIndex(w http.ResponseWriter, r *http.Request) error {
	redirectUrl := os.Getenv("REDIRECT_URL")

	if redirectUrl == "" {
		return nil
	}

	res, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "google",
		RedirectTo: redirectUrl,
	})

	if err != nil {
		return err
	}

	http.Redirect(w, r, res.URL, http.StatusSeeOther)

	return nil
}

func HandlerAuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	err := r.URL.Query().Get("error")

	if err != "" {
		return renderComponent(w, r, auth.ErrorRegister())
	}

	if len(accessToken) == 0 {

		return renderComponent(w, r, auth.CallbackScript())
	}
	setAuthCookie(r, accessToken)
	hxRedirect(w, r, "/")

	return nil
}

func LoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := auth.LoginParams{
		Email: r.FormValue("email"),
	}

	if !validation.IsValidEmail(credentials.Email) {
		credentials.Success = false
		return renderComponent(w, r, auth.LoginForm(credentials, auth.LoginErrors{
			Email: "Invalid email, please try again",
		}))
	}

	err := sb.Client.Auth.SendMagicLink(r.Context(), credentials.Email)

	if err != nil {
		return renderComponent(w, r, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCred: "Invalid credentials, please try again",
		}))
	}

	credentials.Success = true
	return renderComponent(w, r, auth.LoginForm(credentials, auth.LoginErrors{}))
}

func Logout(w http.ResponseWriter, r *http.Request) error {

	session.SessionManager.Destroy(r.Context())
	hxRedirect(w, r, "/login")
	return nil
}

func setAuthCookie(r *http.Request, accessToken string) {
	session.SessionManager.Put(r.Context(), "accessToken", accessToken)
}
