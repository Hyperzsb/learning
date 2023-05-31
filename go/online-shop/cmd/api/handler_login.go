package main

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"onlineshop/internal/model"
	"time"
)

// loginRequest defines the standard request body of the login API,
// which will be sent when the user performs a login operation.
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// loginResponse defines the standard response body of the login API,
// which will be returned in response of the login request.
// The status code states the result of the login operation:
// when it equals "Bad Request", the request body is invalid;
// when it equals "Invalid Credential", the email or password is wrong;
// when it equals "Internal Server Error", some unexpected errors occur.
// In these failure cases, the Token field will be meaningless but not be empty.
type loginResponse struct {
	Status string      `json:"status"`
	Token  model.Token `json:"token"`
}

// login is a handler function to deal with the login operation performed by the user.
// login will return the session token if the credential provided by the user is correct,
// and save the token into database to authenticate and track the user in the future.
func (app *application) login(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	request := loginRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		app.loggers.error.Println(err)
		err = loginBadRequest(w)
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	user, err := app.model.GetUserByEmail(request.Email)
	if err != nil {
		app.loggers.error.Println(err)
		if _, ok := err.(*model.EmptyQueryError); ok {
			err = loginInvalidCredential(w)
			if err != nil {
				app.loggers.error.Println(err)
			}
		} else {
			err = loginInternalServerError(w)
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
			err = loginInvalidCredential(w)
			if err != nil {
				app.loggers.error.Println(err)
			}
		} else {
			err = loginInternalServerError(w)
			if err != nil {
				app.loggers.error.Println(err)
			}
		}

		return
	}

	token, err := model.NewToken(user.ID, "Default Scope", time.Hour*24)
	if err != nil {
		app.loggers.error.Println(err)
		err = loginInternalServerError(w)
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	_, err = app.model.DeleteTokensByUserID(user.ID)
	if err != nil {
		app.loggers.error.Println(err)
		err = loginInternalServerError(w)
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	_, err = app.model.CreateToken(token)
	if err != nil {
		app.loggers.error.Println(err)
		err = loginInternalServerError(w)
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	response := loginResponse{
		Status: "OK",
		Token:  token,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		app.loggers.error.Println(err)
		err = loginInternalServerError(w)
		if err != nil {
			app.loggers.error.Println(err)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseJSON)
}

// loginBadRequest writes the error message to the response
// when the request body does not follow the standard structure.
func loginBadRequest(w http.ResponseWriter) error {
	response := loginResponse{
		Status: "Bad Request",
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(responseJSON)

	return nil
}

// loginInvalidCredential writes the error message to the response
// when the email is not existed or the password is not matched
func loginInvalidCredential(w http.ResponseWriter) error {
	response := loginResponse{
		Status: "Invalid Credential",
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

// loginInternalServerError writes the error message to the response
// when something goes wrong with the internal logics.
func loginInternalServerError(w http.ResponseWriter) error {
	response := loginResponse{
		Status: "Internal Server Error",
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(responseJSON)

	return nil
}
