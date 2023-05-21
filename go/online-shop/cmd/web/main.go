package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	version    = "0.0.1"
	cssVersion = "0.0.1"
)

type configuration struct {
	host   string
	port   int
	env    string
	api    string
	stripe struct {
		key    string
		secret string
	}
}

type application struct {
	version string
	config  configuration
	loggers struct {
		info  *log.Logger
		debug *log.Logger
		error *log.Logger
	}
	templates map[string]*template.Template
}

func (app *application) serve() error {
	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", app.config.host, app.config.port),
		Handler:           app.router(),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	app.loggers.info.Printf("Starting the Web server in %s mode on %s:%d\n", app.config.env, app.config.host, app.config.port)

	return server.ListenAndServe()
}

func main() {
	config := configuration{}

	flag.StringVar(&config.host, "host", "127.0.0.1", "host to listen on")
	flag.IntVar(&config.port, "port", 8080, "port to listen on")
	flag.StringVar(&config.env, "environment", "dev", "serving mode")
	flag.StringVar(&config.api, "api", "127.0.0.1:8000", "url to api")
	flag.Parse()

	config.stripe.key = os.Getenv("STRIPE_KEY")
	config.stripe.secret = os.Getenv("STRIPE_SECRET")

	app := &application{
		version: version,
		config:  config,
		loggers: struct {
			info  *log.Logger
			debug *log.Logger
			error *log.Logger
		}{
			info:  log.New(os.Stdout, "[INFO]  ", log.Ldate|log.Ltime),
			debug: log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
			error: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		},
		templates: make(map[string]*template.Template),
	}

	if err := app.serve(); err != nil {
		app.loggers.error.Fatal(err)
	}
}
