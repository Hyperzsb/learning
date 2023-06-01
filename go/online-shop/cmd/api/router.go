package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *application) router() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Post("/login", app.login)
	mux.Post("/authenticate", app.authenticate)

	mux.Route("/product", func(r chi.Router) {
		r.Post("/", app.createProduct)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", app.getProduct)
			r.Put("/", app.updateProduct)
			r.Delete("/", app.deleteProduct)
		})
	})

	mux.Route("/admin", func(r chi.Router) {
		r.Use(app.needAuthentication)

		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Success"))
		})
	})

	mux.Post("/payment", app.createPaymentIntent)

	return mux
}
