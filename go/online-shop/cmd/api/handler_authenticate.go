package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"onlineshop/internal/model"
	"strings"
)

// authenticateRequest defines the standard request body of the authenticate API,
// which will be sent when the user checks the validity of a token. Currently,
// there will be no data sent along with the authenticateRequest.
//
type authenticateRequest struct {
}

// authenticateResponse defines the standard response body of the login API,
// which will be returned in response of the login request.
type authenticateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// authenticate checks whether the current user has logged in and has a valid token.
func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	user, err := app.authenticateToken(r)
	if err != nil {
		app.loggers.error.Println(err)
		err = authenticateInvalidToken(w)
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	err = authenticateValidToken(w)
	if err != nil {
		app.loggers.error.Println(err)
		return
	}

	_ = user
}

// authenticateToken reads the request and retrieves the token from the
// Authorization header if exists. After getting the token, it will retrieve
// the corresponding user of this token.
func (app *application) authenticateToken(r *http.Request) (model.User, error) {
	user := model.User{}

	authorizationHeader := r.Header.Get("Authorization")
	if len(authorizationHeader) == 0 {
		return user, errors.New("no authorization header provided")
	}

	authorizationParts := strings.Split(authorizationHeader, " ")

	if len(authorizationParts) != 2 ||
		authorizationParts[0] != "Bearer" {
		return user, errors.New("invalid authorization header")
	}

	user, err := app.model.GetUserByToken(authorizationParts[1])
	if err != nil {
		app.loggers.error.Println(err)
		if _, ok := err.(*model.EmptyQueryError); ok {
			return user, errors.New("invalid authorization token")
		}

		return user, errors.New("internal server error")
	}

	return user, nil
}

// authenticateInvalidToken writes the error message to the response
// when the token is not provided or the token is invalid.
func authenticateInvalidToken(w http.ResponseWriter) error {
	response := authenticateResponse{
		Status: "Invalid Token",
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write(responseJSON)

	return nil
}

// authenticateInvalidToken writes the success message to the response
// when the given token is invalid.
func authenticateValidToken(w http.ResponseWriter) error {
	response := authenticateResponse{
		Status: "Valid Token",
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseJSON)

	return nil
}
