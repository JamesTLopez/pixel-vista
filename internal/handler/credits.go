package handler

import (
	"net/http"
	"os"
	"pixelvista/view/pages/credits"

	"github.com/stripe/stripe-go/v79"
)

func HandleCreditsIndex(w http.ResponseWriter, r *http.Request) error {

	return renderComponent(w, r, credits.CreditsIndex())
}

func StripeCallbackSuccess(w http.ResponseWriter, r *http.Request) error {
	stripe.Key = os.Getenv("STRIPE_API_KEY")
	return nil
}
