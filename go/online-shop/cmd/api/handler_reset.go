package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"net/http"
	"onlineshop/cmd/api/jsonio"
	"time"
)

//go:embed templates
var templateFS embed.FS

func (app *application) reset(w http.ResponseWriter, r *http.Request) {
	type resetRequest struct {
		Email string `json:"email"`
	}

	request := resetRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		app.loggers.error.Println(err)
		err = jsonio.Write(w, jsonio.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Your request format is invalid. Please try again.",
		})
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	body := struct {
		Link string
	}{"https://google.com"}

	err = app.sendMail("from@onlineshop.com", "to@onlineshop.com", "Reset Password", "reset", body)
	if err != nil {
		app.loggers.error.Println(err)
		app.loggers.error.Println(err)
		err = jsonio.Write(w, jsonio.Response{
			Code:    http.StatusInternalServerError,
			Status:  "Interval Server Error",
			Message: "Unexpected errors occurred. Please try again later.",
		})
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	err = jsonio.Write(w, jsonio.Response{
		Code:    http.StatusOK,
		Status:  "Reset Email Sent",
		Message: "The reset email has been sent. Please check your inbox",
	})
	if err != nil {
		app.loggers.error.Println(err)
	}
}

func (app *application) sendMail(from, to, subject, tmpl string, body interface{}) error {
	htmlTemplate, err := template.New(fmt.Sprintf("%s.html.tmpl", tmpl)).ParseFS(templateFS, fmt.Sprintf("templates/%s.html.gohtml", tmpl))
	if err != nil {
		return err
	}

	htmlTemplateData := bytes.Buffer{}
	err = htmlTemplate.ExecuteTemplate(&htmlTemplateData, "body", body)
	if err != nil {
		return err
	}

	htmlMessage := htmlTemplateData.String()

	plainTemplate, err := template.New(fmt.Sprintf("%s.plain.tmpl", tmpl)).ParseFS(templateFS, fmt.Sprintf("templates/%s.plain.gohtml", tmpl))
	if err != nil {
		return err
	}

	plainTemplateData := bytes.Buffer{}
	err = plainTemplate.ExecuteTemplate(&plainTemplateData, "body", body)
	if err != nil {
		return err
	}

	plainMessage := plainTemplateData.String()

	server := mail.NewSMTPClient()
	server.Host = app.config.mailtrap.smtp.host
	server.Port = app.config.mailtrap.smtp.port
	server.Username = app.config.mailtrap.smtp.username
	server.Password = app.config.mailtrap.smtp.password
	server.Encryption = mail.EncryptionTLS
	server.KeepAlive = false
	server.ConnectTimeout = time.Second * 10
	server.SendTimeout = time.Second * 10

	client, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from).AddTo(to).SetSubject(subject)
	email.SetBody(mail.TextHTML, htmlMessage)
	email.AddAlternative(mail.TextPlain, plainMessage)

	err = email.Send(client)
	if err != nil {
		return err
	}

	return nil
}
