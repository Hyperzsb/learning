package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) route() http.Handler {
	mux := chi.NewRouter()
	mux.Get("/terminal", app.terminal)
	mux.Get("/about", app.about)

	return mux
}
