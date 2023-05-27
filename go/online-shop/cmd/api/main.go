package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"onlineshop/internal/model"
	"os"
	"time"
)

const (
	version = "0.0.1"
)

type configuration struct {
	host string
	port int
	env  string
	db   struct {
		dsn string
	}
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
	model *model.Model
}

func (app *application) serve() error {
	var err error
	app.model, err = model.New(app.config.db.dsn)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", app.config.host, app.config.port),
		Handler:           app.router(),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	app.loggers.info.Printf("Starting the API server in %s mode on %s:%d\n", app.config.env, app.config.host, app.config.port)

	return server.ListenAndServe()
}

func main() {
	config := configuration{}

	flag.StringVar(&config.host, "host", "127.0.0.1", "host to listen on")
	flag.IntVar(&config.port, "port", 8000, "port to listen on")
	flag.StringVar(&config.env, "environment", "dev", "serving mode")
	flag.Parse()

	config.db.dsn = os.Getenv("API_DSN")
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
			info:  log.New(os.Stdout, "[API|INFO]  ", log.Ldate|log.Ltime),
			debug: log.New(os.Stdout, "[API|DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
			error: log.New(os.Stderr, "[API|ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		},
	}

	if err := app.serve(); err != nil {
		_ = app.model.Close()
		app.loggers.error.Fatal(err)
	}
}
