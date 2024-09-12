package handler

import (
	"net/http"
	"os"
	"pixelvista/db"
	"pixelvista/types"
	"pixelvista/view/pages/credits"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
)

func HandleCreditsIndex(w http.ResponseWriter, r *http.Request) error {

	prices, err := db.GetCreditPrices()

	if err != nil {
		return renderComponent(w, r, credits.CreditsIndex([]types.CreditPrice{}))
	}

	return renderComponent(w, r, credits.CreditsIndex(prices))
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
