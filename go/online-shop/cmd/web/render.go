package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

type templateData struct {
	version    string
	cssVersion string
}

//go:embed templates
var templateFS embed.FS

func (app *application) render(w http.ResponseWriter, r *http.Request, page string, td *templateData) error {
	if td == nil {
		td = &templateData{}
	}

	var tmpl *template.Template
	var err error

	if t, ok := app.templates[page]; !ok {
		tmpl, err = app.parse(page)
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

func (app *application) parse(page string) (*template.Template, error) {
	var tmpl *template.Template
	var err error

	tmpl, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).ParseFS(templateFS, "templates/base.layout.gohtml", fmt.Sprintf("templates/%s.page.gohtml", page))
	if err != nil {
		app.loggers.error.Println(err)
		return nil, err
	}

	app.templates[page] = tmpl

	return tmpl, nil
}
