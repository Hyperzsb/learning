package payment

import (
	"encoding/json"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

type Payment struct {
	key    string
	secret string
}

func New(key, secret string) *Payment {
	return &Payment{key: key, secret: secret}
}

func (p *Payment) Charge(currency string, amount int) ([]byte, error) {
	return p.chargeByStripe(currency, amount)
}

func (p *Payment) chargeByStripe(currency string, amount int) ([]byte, error) {
	stripe.Key = p.secret
	params := &stripe.PaymentIntentParams{
		Currency: stripe.String(currency),
		Amount:   stripe.Int64(int64(amount)),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok {
			err = &ChargeError{stripeErr.Error()}
		}

		return nil, err
	}

	data, err := json.Marshal(pi)
	if err != nil {
		return nil, err
	}

	return data, nil
}
