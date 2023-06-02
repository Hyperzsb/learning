package main

import "net/http"

func (app *application) needSession(next http.Handler) http.Handler {
	return app.session.LoadAndSave(next)
}

// needAuthentication is a middleware that checks the authentication status
// of incoming requests. If requests are sent without an Authorization header
// or the token inside the header is not valid, it will return an error
// indicating that the token is not valid.
func (app *application) needAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := app.authenticateToken(r)
		if err != nil {
			app.loggers.error.Println(err)
			err = authenticateInvalidToken(w)
			if err != nil {
				app.loggers.error.Println(err)
			}

			return
		}

		next.ServeHTTP(w, r)
	})
}
