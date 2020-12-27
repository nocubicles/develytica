package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nocubicles/develytica/src/models"
	"github.com/nocubicles/develytica/src/utils"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/checkout/session"
)

type BillingPageData struct {
	Authenticated    bool
	UserName         string
	ValidationErrors map[string]string
	SessionID        string
}

func HandleBillingSetup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		PublishableKey string `json:"publishableKey"`
		SmallPrice     string `json:"smallPrice"`
		MediumPrice    string `json:"mediumPrice"`
		LargePrice     string `json:"largePrice"`
	}{
		PublishableKey: os.Getenv("STRIPE_PUBLIC"),
		SmallPrice:     fmt.Sprint("price_1I2hfsEuUDslru3nmuQ2VOvW"),
		MediumPrice:    fmt.Sprint("price_1I2hgHEuUDslru3nOwoYrthP"),
		LargePrice:     fmt.Sprint("price_1I2hgUEuUDslru3nLzgvPntG"),
	})
}

func HandleCheckoutCallback(w http.ResponseWriter, r *http.Request) {

	sessionID := r.URL.Query().Get("session_id")

	s, err := session.Get(sessionID, nil)
	if err != nil {
		fmt.Println(err)
	}
	db := utils.DbConnection()

	db.Model(&models.Tenant{}).Where("id = ?",
		convertStringToUint(s.ClientReferenceID)).
		Update("stripe_id", s.Customer.ID)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func HandleCreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	billingSecret := os.Getenv("STRIPE_SECRET")
	stripe.Key = billingSecret
	server := os.Getenv("SERVER")
	successURL := server + "/billing/success?session_id={CHECKOUT_SESSION_ID}"
	cancelURL := server + "/billing/fail"
	authContext, _ := getAuthContextData(r)

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodPost {

		var req struct {
			PriceID string `json:"priceID"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("json.NewDecoder.Decode: %v", err)
			return
		}

		params := &stripe.CheckoutSessionParams{
			SuccessURL:         &successURL,
			CancelURL:          &cancelURL,
			ClientReferenceID:  stripe.String(fmt.Sprint(authContext.TenantID)),
			PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
			Mode:               stripe.String(string(stripe.CheckoutSessionModeSubscription)),
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String(req.PriceID),
					Quantity: stripe.Int64(1),
				}},
		}

		s, err := session.New(params)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(struct {
				ErrorData string `json:"error"`
			}{
				ErrorData: err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			SessionID string `json:"sessionID"`
		}{
			SessionID: s.ID,
		})

	}
}
