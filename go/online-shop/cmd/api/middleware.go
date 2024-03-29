package main

import (
	"net/http"
	"onlineshop/cmd/api/jsonio"
)

func (app *application) needsSession(next http.Handler) http.Handler {
	return app.session.LoadAndSave(next)
}

// needsAuthorization is a middleware that checks the authorization status of
// incoming requests. If requests are sent without an Authorization header
// or the token inside the header is not valid, it will return an error
// indicating that the token is invalid.
func (app *application) needsAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := app.validateToken(r)
		if err != nil {
			app.loggers.error.Println(err)
			err = jsonio.Write(w, jsonio.Response{
				Code:    http.StatusForbidden,
				Status:  "Invalid Token",
				Message: "No valid token is provided. Please login first.",
			})
			if err != nil {
				app.loggers.error.Println(err)
			}

			return
		}

		next.ServeHTTP(w, r)
	})
}
