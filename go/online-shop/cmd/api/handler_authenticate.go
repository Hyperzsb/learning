package main

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"onlineshop/cmd/api/jsonio"
	"onlineshop/internal/model"
	"time"
)

// authenticateRequest defines the standard request body of the authorize API,
// which will be sent when the user performs a authenticate operation.
type authenticateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// authenticate is a handler function to deal with the authenticate operation
// performed by the user. authenticate will return the session token if the
// credential provided by the user is correct, and save the token into database
// to authorize and track the user in the future.
func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	request := authenticateRequest{}
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

	user, err := app.model.GetUserByEmail(request.Email)
	if err != nil {
		app.loggers.error.Println(err)
		if _, ok := err.(*model.EmptyQueryError); ok {
			err = jsonio.Write(w, jsonio.Response{
				Code:    http.StatusUnauthorized,
				Status:  "Invalid Credential",
				Message: "Your credential (email or password) is invalid. Please try again.",
			})
			if err != nil {
				app.loggers.error.Println(err)
			}
		} else {
			err = jsonio.Write(w, jsonio.Response{
				Code:    http.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Some unexpected errors occur. Please try again later.",
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
			err = jsonio.Write(w, jsonio.Response{
				Code:    http.StatusUnauthorized,
				Status:  "Invalid Credential",
				Message: "Your credential (email or password) is invalid. Please try again.",
			})
			if err != nil {
				app.loggers.error.Println(err)
			}
		} else {
			err = jsonio.Write(w, jsonio.Response{
				Code:    http.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Some unexpected errors occur. Please try again later.",
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
		err = jsonio.Write(w, jsonio.Response{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Some unexpected errors occur. Please try again later.",
		})
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	_, err = app.model.DeleteTokensByUserID(user.ID)
	if err != nil {
		app.loggers.error.Println(err)
		err = jsonio.Write(w, jsonio.Response{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Some unexpected errors occur. Please try again later.",
		})
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	_, err = app.model.CreateToken(token)
	if err != nil {
		app.loggers.error.Println(err)
		err = jsonio.Write(w, jsonio.Response{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Some unexpected errors occur. Please try again later.",
		})
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	err = jsonio.Write(w, jsonio.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Your credential is valid. Token has been generated.",
		Data:    token,
	})
	if err != nil {
		app.loggers.error.Println(err)
	}
}
