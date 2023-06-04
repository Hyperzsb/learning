package main

import (
	"errors"
	"fmt"
	"net/http"
	"onlineshop/internal/model"
	"strings"
)

// authorizeRequest defines the standard request body of the authorization API,
// which will be sent when the user checks the validity of a token. Currently,
// there will be no data sent along with the authorizeRequest, so this type
// definition is left blank intentionally.
type authorizeRequest struct {
}

// authorizeResponse defines the standard response body of the authenticate API,
// which will be returned in response of the authenticate request.
// It implements the GeneralResponse interface.
// The status states the result of the authorization operation:
// when it equals to "Invalid Token", the token is invalid or expired, or some
// unexpected interval errors occur;
// when it equals to "Valid Token", the token is valid and the authorization is approved.
// The message of the response will indicate the detailed reason of failure.
type authorizeResponse struct {
	code    int
	status  string
	message string
}

func (ar authorizeResponse) Code() int {
	return ar.code
}

func (ar authorizeResponse) Status() string {
	return ar.status
}

func (ar authorizeResponse) Message() string {
	return ar.message
}

// authorize checks whether the current user has logged in and has a valid token.
func (app *application) authorize(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	user, err := app.validateToken(r)
	if err != nil {
		app.loggers.error.Println(err)
		err = writeJSON(w, authorizeResponse{
			code:    http.StatusForbidden,
			status:  "Invalid Token",
			message: fmt.Sprintf("Token check failed (%s). Please log in again.", err.Error()),
		})
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	err = writeJSON(w, authorizeResponse{
		code:    http.StatusOK,
		status:  "Valid Token",
		message: "Your token is valid. Request is permitted.",
	})
	if err != nil {
		app.loggers.error.Println(err)
		return
	}

	_ = user
}

// validateToken reads the request and retrieves the token from the
// Authorization header if exists. After getting the token, it will retrieve
// the corresponding user of this token.
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
