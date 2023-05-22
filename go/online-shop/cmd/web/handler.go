package main

import "net/http"

func (app *application) checkout(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	if err := app.render(w, r, "checkout", nil); err != nil {
		return
	}
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	if err := app.render(w, r, "about", nil); err != nil {
		return
	}
}
