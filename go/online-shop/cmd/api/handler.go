package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"onlineshop/internal/model"
	"onlineshop/internal/payment"
	"strconv"
)

func (app *application) createPaymentIntent(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	type paymentRequest struct {
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	}

	request := paymentRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	amount, err := strconv.Atoi(request.Amount)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	p := payment.New(app.config.stripe.key, app.config.stripe.secret)

	data, err := p.Charge(request.Currency, amount)
	if err != nil {
		app.loggers.error.Println(err)
		if chargeErr, ok := err.(*payment.ChargeError); ok {
			http.Error(w, chargeErr.Error(), http.StatusForbidden)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(data)
}

func (app *application) createProduct(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	type productResponse struct {
		ID int `json:"id"`
	}

	product := model.Product{}
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := app.model.CreateProduct(product)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := productResponse{
		ID: id,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(responseJSON)
}

func (app *application) getProduct(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	rawID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	product, err := app.model.GetProduct(id)
	if err != nil {
		app.loggers.error.Println(err)
		if _, ok := err.(*model.EmptyQueryError); ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		return
	}

	productJSON, err := json.Marshal(product)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(productJSON)
}

func (app *application) updateProduct(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	type productResponse struct {
		CNT int `json:"cnt"`
	}

	rawID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	product := model.Product{}
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	cnt, err := app.model.UpdateProduct(id, product)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := productResponse{
		CNT: cnt,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseJSON)
}

func (app *application) deleteProduct(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	type productResponse struct {
		CNT int `json:"cnt"`
	}

	rawID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	cnt, err := app.model.DeleteProduct(id)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := productResponse{
		CNT: cnt,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseJSON)
}
