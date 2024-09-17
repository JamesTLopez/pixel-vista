package handler

import (
	"fmt"
	"net/http"
	"os"
	"pixelvista/db"
	"pixelvista/internal"
	"pixelvista/types"
	"pixelvista/view/pages/credits"

	"github.com/go-chi/chi/v5"
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
	productID := chi.URLParam(r, "productID")
	stripe.Key = os.Getenv("STRIPE_API_KEY")

	checkoutParams := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(os.Getenv("STRIPE_SUCCESS_URL")),
		CancelURL:  stripe.String(os.Getenv("STRIPE_CANCEL_URL")),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(productID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}

	s, err := session.New(checkoutParams)

	if err != nil {
		return err
	}
	return hxRedirect(w, r, s.URL)
}

func StripeCheckoutSuccess(w http.ResponseWriter, r *http.Request) error {
	user := internal.GetAuthenticatedUser(r)
	sessionID := chi.URLParam(r, "sessionID")

	_, err := session.Get(sessionID, nil)
	if err != nil {
		fmt.Println("prevent stripesession charge cause it failed")
		return err
	}

	lineItems := session.ListLineItems(&stripe.CheckoutSessionListLineItemsParams{
		Session: stripe.String(sessionID),
	})

	lineItems.Next()
	currentItem := lineItems.LineItem()

	priceID := currentItem.Price.ID

	credits, err := db.GetCreditPriceByID(priceID)

	if err != nil {
		fmt.Println("prevent stripesession charge cause it failed")
		return err
	}
	fmt.Println(priceID, credits.Credits)

	if credits.Credits <= 0 {
		fmt.Println("prevent stripesession charge cause it failed")
		return nil
	}

	user.Account.Credits = credits.Credits

	if err := db.UpdateProfile(&user.Account); err != nil {
		fmt.Println("test------", err)
		return err
	}

	hxRedirect(w, r, "/generate")

	return nil
}

func StripeCheckoutCancel(w http.ResponseWriter, r *http.Request) error {
	return nil
}
