package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) router() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/checkout", app.checkout)
	mux.Post("/receipt", app.receipt)
	mux.Get("/about", app.about)

	fs := http.FileServer(http.Dir("./cmd/web/static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fs))

	return mux
}
