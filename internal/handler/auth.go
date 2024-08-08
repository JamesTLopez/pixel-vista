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
func HandlerRegisterIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.Register().Render(r.Context(), w)
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
	setAuthCookie(r, resp.AccessToken)
	hxRedirect(w, r, "/")
	return nil
}

func RegisterCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.RegisterParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}

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

func Logout(w http.ResponseWriter, r *http.Request) error {

	session.SessionManager.Destroy(r.Context())
	hxRedirect(w, r, "/login")
	return nil
}

func setAuthCookie(r *http.Request, accessToken string) {
	session.SessionManager.Put(r.Context(), "accessToken", accessToken)
}
