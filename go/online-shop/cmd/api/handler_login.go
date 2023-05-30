package main

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"onlineshop/internal/model"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	app.loggers.info.Printf("%s -> %s\n", r.Method, r.URL)

	request := loginRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, "{}", http.StatusBadRequest)
		return
	}

	user, err := app.model.GetUserByEmail(request.Email)
	if err != nil {
		app.loggers.error.Println(err)
		if _, ok := err.(*model.EmptyQueryError); ok {
			err = invalidCredential(w)
			if err != nil {
				app.loggers.error.Println(err)
			}
		} else {
			http.Error(w, "{}", http.StatusInternalServerError)
		}

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		app.loggers.error.Println(err)
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			err = invalidCredential(w)
			if err != nil {
				app.loggers.error.Println(err)
			}
		} else {
			http.Error(w, "{}", http.StatusInternalServerError)
		}

		return
	}

	response := loginResponse{
		Status: "OK",
		Token:  "Token",
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		app.loggers.error.Println(err)
		http.Error(w, "{}", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseJSON)
}

// invalidCredential writes the error message to the response when
// the email is not existed or the password is not matched
func invalidCredential(w http.ResponseWriter) error {
	response := loginResponse{
		Status: "Invalid Credential",
		Token:  "",
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
