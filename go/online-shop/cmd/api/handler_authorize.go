package main

import (
	"errors"
	"fmt"
	"net/http"
	"onlineshop/cmd/api/jsonio"
	"onlineshop/internal/model"
	"strings"
)

// authorize checks whether the current user has logged in and has a valid token.
func (app *application) authorize(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	user, err := app.validateToken(r)
	if err != nil {
		app.loggers.error.Println(err)
		err = jsonio.Write(w, jsonio.Response{
			Code:    http.StatusForbidden,
			Status:  "Invalid Token",
			Message: fmt.Sprintf("Token check failed (%s). Please log in again.", err.Error()),
		})
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	err = jsonio.Write(w, jsonio.Response{
		Code:    http.StatusOK,
		Status:  "Valid Token",
		Message: "Your token is valid. Request is permitted.",
	})
	if err != nil {
		app.loggers.error.Println(err)
		return
	}

	_ = user
}

// validateToken reads the request and retrieves the token from the Authorization
// header if exists. After getting the token, it will retrieve the corresponding
// user of this token.
func (app *application) validateToken(r *http.Request) (model.User, error) {
	user := model.User{}

	authorizationHeader := r.Header.Get("Authorization")
	if len(authorizationHeader) == 0 {
		return user, errors.New("no authorization header provided")
	}

	authorizationParts := strings.Split(authorizationHeader, " ")

	if len(authorizationParts) != 2 ||
		authorizationParts[0] != "Bearer" {
		return user, errors.New("invalid authorization header, not using Bearer scheme")
	}

	user, err := app.model.GetUserByToken(authorizationParts[1])
	if err != nil {
		app.loggers.error.Println(err)
		if _, ok := err.(*model.EmptyQueryError); ok {
			return user, errors.New("invalid or expired authorization token")
		}

		return user, errors.New("internal server error")
	}

	return user, nil
}
