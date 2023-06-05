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

	mux.Use(app.needsSession)

	mux.Post("/authenticate", app.authenticate)
	mux.Post("/deauthenticate", app.deauthenticate)
	mux.Post("/authorize", app.authorize)

	mux.Post("/reset", app.reset)

	mux.Route("/product", func(r chi.Router) {
		r.Post("/", app.createProduct)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", app.getProduct)
			r.Put("/", app.updateProduct)
			r.Delete("/", app.deleteProduct)
		})
	})

	mux.Post("/payment", app.createPaymentIntent)

	mux.Route("/admin", func(r chi.Router) {
		r.Use(app.needsAuthorization)

		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Success"))
		})
	})

	// This group of routes is used for testing purposes only.
	mux.Route("/test", func(r chi.Router) {
		r.Route("/session", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				err := r.ParseForm()
				if err != nil {
					_, _ = w.Write([]byte("Post Failed"))
				}

				data := r.Form.Get("data")
				app.session.Put(r.Context(), "data", data)
				_, _ = w.Write([]byte("Post Succeeded: " + data))
			})

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				data := app.session.GetString(r.Context(), "data")
				_, _ = w.Write([]byte("Get Succeeded: " + data))
			})

			r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				err := app.session.Destroy(r.Context())
				if err != nil {
					_, _ = w.Write([]byte("Delete Failed"))
				}

				_, _ = w.Write([]byte("Delete Succeeded"))
			})
		})
	})

	return mux
}
