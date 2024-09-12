package handler

import (
	"net/http"
	"os"
	"pixelvista/view/pages/credits"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
)

func HandleCreditsIndex(w http.ResponseWriter, r *http.Request) error {

	return renderComponent(w, r, credits.CreditsIndex())
}

func StripeCheckout(w http.ResponseWriter, r *http.Request) error {
	stripe.Key = os.Getenv("STRIPE_API_KEY")
	checkoutParams := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(""),
		CancelURL:  stripe.String(""),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String("__"),
				Quantity: stripe.Int64(1),
			},
		},
	}

	s, err := session.New(checkoutParams)

	if err != nil {
		return err
	}
	http.Redirect(w, r, s.URL, http.StatusSeeOther)
	return nil
}
