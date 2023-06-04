package main

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"onlineshop/internal/model"
	"time"
)

// authenticateRequest defines the standard request body of the authorize API,
// which will be sent when the user performs a authenticate operation.
type authenticateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// authenticateResponse defines the standard response body of the authorize API,
// which will be returned in response of the authenticate request.
// It implements the GeneralResponse interface.
// The status states the result of the authentication operation:
// when it equals to "Bad Request", the request body is invalid;
// when it equals to "Invalid Credential", the email or password is wrong;
// when it equals to "Internal Server Error", some unexpected errors occur.
// In these failure cases, the Token field will be meaningless but not be empty.
type authenticateResponse struct {
	code    int
	status  string
	message string
	Token   model.Token `json:"token"`
}

func (ar authenticateResponse) Code() int {
	return ar.code
}

func (ar authenticateResponse) Status() string {
	return ar.status
}

func (ar authenticateResponse) Message() string {
	return ar.message
}

// authenticate is a handler function to deal with the authenticate operation performed by the user.
// authenticate will return the session token if the credential provided by the user is correct,
// and save the token into database to authorize and track the user in the future.
func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	request := authenticateRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		app.loggers.error.Println(err)
		err = writeJSON(w, authenticateResponse{
			code:    http.StatusBadRequest,
			status:  "Bad Request",
			message: "Your request format is invalid. Please try again.",
		})
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	user, err := app.model.GetUserByEmail(request.Email)
	if err != nil {
		app.loggers.error.Println(err)
		if _, ok := err.(*model.EmptyQueryError); ok {
			err = writeJSON(w, authenticateResponse{
				code:    http.StatusUnauthorized,
				status:  "Invalid Credential",
				message: "Your credential (email or password) is invalid. Please try again.",
			})
			if err != nil {
				app.loggers.error.Println(err)
			}
		} else {
			err = writeJSON(w, authenticateResponse{
				code:    http.StatusInternalServerError,
				status:  "Internal Server Error",
				message: "Some unexpected errors occur. Please try again later.",
			})
			if err != nil {
				app.loggers.error.Println(err)
			}
		}

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		app.loggers.error.Println(err)
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			err = writeJSON(w, authenticateResponse{
				code:    http.StatusUnauthorized,
				status:  "Invalid Credential",
				message: "Your credential (email or password) is invalid. Please try again.",
			})
			if err != nil {
				app.loggers.error.Println(err)
			}
		} else {
			err = writeJSON(w, authenticateResponse{
				code:    http.StatusInternalServerError,
				status:  "Internal Server Error",
				message: "Some unexpected errors occur. Please try again later.",
			})
			if err != nil {
				app.loggers.error.Println(err)
			}
		}

		return
	}

	token, err := model.NewToken(user.ID, "Default Scope", time.Hour*24)
	if err != nil {
		app.loggers.error.Println(err)
		err = writeJSON(w, authenticateResponse{
			code:    http.StatusInternalServerError,
			status:  "Internal Server Error",
			message: "Some unexpected errors occur. Please try again later.",
		})
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	_, err = app.model.DeleteTokensByUserID(user.ID)
	if err != nil {
		app.loggers.error.Println(err)
		err = writeJSON(w, authenticateResponse{
			code:    http.StatusInternalServerError,
			status:  "Internal Server Error",
			message: "Some unexpected errors occur. Please try again later.",
		})
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	_, err = app.model.CreateToken(token)
	if err != nil {
		app.loggers.error.Println(err)
		err = writeJSON(w, authenticateResponse{
			code:    http.StatusInternalServerError,
			status:  "Internal Server Error",
			message: "Some unexpected errors occur. Please try again later.",
		})
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	err = writeJSON(w, authenticateResponse{
		code:    http.StatusOK,
		status:  "OK",
		message: "Your credential is valid. Token has been generated.",
		Token:   token,
	})
	if err != nil {
		app.loggers.error.Println(err)
	}
}
