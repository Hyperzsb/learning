package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

type templateData struct {
	Version    string
	CSSVersion string
	API        string
	Funcs      template.FuncMap
	Data       map[string]interface{}
}

func (app *application) initDefaultTemplateData(td *templateData) *templateData {
	if td == nil {
		td = &templateData{
			Version:    version,
			CSSVersion: cssVersion,
			API:        app.config.api,
			Funcs:      template.FuncMap{},
			Data:       make(map[string]interface{}),
		}

		td.Data["StripeKey"] = app.config.stripe.key

		return td
	}

	if td.Version == "" {
		td.Version = version
	}

	if td.CSSVersion == "" {
		td.CSSVersion = cssVersion
	}

	if td.API == "" {
		td.API = app.config.api
	}

	if td.Funcs == nil {
		td.Funcs = template.FuncMap{}
	}

	if td.Data == nil {
		td.Data = make(map[string]interface{})
	}

	if _, ok := td.Data["StripeKey"]; !ok {
		td.Data["StripeKey"] = app.config.stripe.key
	}

	return td
}

//go:embed templates
var templateFS embed.FS

func (app *application) render(w http.ResponseWriter, r *http.Request, page string, td *templateData) error {
	td = app.initDefaultTemplateData(td)

	var tmpl *template.Template
	var err error

	if t, ok := app.templates[page]; !ok {
		tmpl, err = app.parse(page, td.Funcs)
		if err != nil {
			app.loggers.error.Println(err)
			return err
		}
	} else {
		tmpl = t
	}

	if err = tmpl.Execute(w, td); err != nil {
		app.loggers.error.Println(err)
		return err
	}

	return nil
}

func (app *application) parse(page string, funcs template.FuncMap) (*template.Template, error) {
	var tmpl *template.Template
	var err error

	tmpl, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(funcs).ParseFS(templateFS, "templates/base.layout.gohtml", fmt.Sprintf("templates/%s.page.gohtml", page))
	if err != nil {
		app.loggers.error.Println(err)
		return nil, err
	}

	app.templates[page] = tmpl

	return tmpl, nil
}
