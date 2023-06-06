package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"onlineshop/internal/model"
	"os"
	"strconv"
	"time"
)

const (
	version = "0.0.1"
)

type configuration struct {
	host   string
	port   int
	env    string
	crypto struct {
		rsa struct {
			pk *rsa.PublicKey
			sk *rsa.PrivateKey
		}
	}
	db struct {
		dsn string
	}
	stripe struct {
		key    string
		secret string
	}
	mailtrap struct {
		smtp struct {
			host     string
			port     int
			username string
			password string
		}
	}
}

func (cfg *configuration) init() error {
	flag.StringVar(&cfg.host, "host", "127.0.0.1", "host to listen on")
	flag.IntVar(&cfg.port, "port", 8000, "port to listen on")
	flag.StringVar(&cfg.env, "environment", "dev", "serving mode")
	flag.Parse()

	// Inject the Stripe config via environment variables
	cfg.db.dsn = os.Getenv("API_DSN")
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	// Inject the Mailtrap config via environment variables
	cfg.mailtrap.smtp.host = os.Getenv("MAILTRAP_SMTP_HOST")
	cfg.mailtrap.smtp.port, _ = strconv.Atoi(os.Getenv("MAILTRAP_SMTP_PORT"))
	cfg.mailtrap.smtp.username = os.Getenv("MAILTRAP_SMTP_USERNAME")
	cfg.mailtrap.smtp.password = os.Getenv("MAILTRAP_SMTP_PASSWORD")

	if err := cfg.loadCrypto(); err != nil {
		return err
	}

	return nil
}

func (cfg *configuration) loadCrypto() error {
	// Inject the file path of crypto keys via environment variables
	rsaPKPath := os.Getenv("RSA_PUBLIC_KEY_PATH")
	if rsaPKPath == "" {
		rsaPKPath = ".config/rsa.pub"
	}
	rsaSKPath := os.Getenv("RSA_PRIVATE_KEY_PATH")
	if rsaSKPath == "" {
		rsaSKPath = ".config/rsa.pem"
	}

	rsaPKBytes, err := os.ReadFile(rsaPKPath)
	if err != nil {
		return err
	}

	rsaSKBytes, err := os.ReadFile(rsaSKPath)
	if err != nil {
		return err
	}

	rsaPKBlock, _ := pem.Decode(rsaPKBytes)
	rsaSKBlock, _ := pem.Decode(rsaSKBytes)

	rsaPK, err := x509.ParsePKCS1PublicKey(rsaPKBlock.Bytes)
	if err != nil {
		return err
	}

	rsaSK, err := x509.ParsePKCS1PrivateKey(rsaSKBlock.Bytes)
	if err != nil {
		return err
	}

	cfg.crypto.rsa.pk = rsaPK
	cfg.crypto.rsa.sk = rsaSK

	return nil
}

type application struct {
	version string
	config  *configuration
	loggers struct {
		info  *log.Logger
		debug *log.Logger
		error *log.Logger
	}
	model   *model.Model
	session *scs.SessionManager
}

func (app *application) serve() error {
	var err error
	app.model, err = model.New(app.config.db.dsn)
	if err != nil {
		return err
	}

	app.session = scs.New()
	app.session.Lifetime = time.Hour * 24
	app.session.Store = mysqlstore.New(app.model.DB())

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

	app := &application{
		version: version,
		config:  &config,
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

	if err := config.init(); err != nil {
		app.loggers.error.Fatal(err)
	}

	if err := app.serve(); err != nil {
		_ = app.model.Close()
		app.loggers.error.Fatal(err)
	}
}
