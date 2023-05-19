package main

import (
	"encoding/json"
	"net/http"
	"onlineshop/internal/payment"
	"strconv"
)

type paymentPayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

func (app *application) createPaymentIntent(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	payload := paymentPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	p := payment.New(app.config.stripe.key, app.config.stripe.secret)

	data, err := p.Charge(payload.Currency, amount)
	if err != nil {
		app.loggers.error.Println(err)
		if chargeErr, ok := err.(*payment.ChargeError); ok {
			http.Error(w, chargeErr.Error(), http.StatusForbidden)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
