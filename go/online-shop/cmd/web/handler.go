package main

import "net/http"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	if err := app.render(w, r, "home", nil); err != nil {
		app.loggers.error.Println(err)
		return
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	if err := app.render(w, r, "login", nil); err != nil {
		app.loggers.error.Println(err)
		return
	}
}

func (app *application) checkout(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	if err := app.render(w, r, "checkout", nil); err != nil {
		app.loggers.error.Println(err)
		return
	}
}

func (app *application) receipt(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	err := r.ParseForm()
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data := make(map[string]interface{})
	data["ID"] = r.Form.Get("payment-id")
	data["Method"] = r.Form.Get("payment-method")
	data["Currency"] = r.Form.Get("payment-currency")
	data["Amount"] = r.Form.Get("payment-amount")
	data["Email"] = r.Form.Get("payment-email")

	if err := app.render(w, r, "receipt", &templateData{Data: data}); err != nil {
		app.loggers.error.Println(err)
		return
	}
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	if err := app.render(w, r, "about", nil); err != nil {
		app.loggers.error.Println(err)
		return
	}
}
