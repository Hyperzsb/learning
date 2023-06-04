package main

import (
	"encoding/json"
	"net/http"
)

type GeneralRequest interface {
}

type GeneralResponse interface {
	Code() int
	Status() string
	Message() string
}

// writeJSON is a utility function to write arbitrary data to the JSON response.
// The response struct always obeys a standard struct providing essential status
// and message to the client. It may fail to write a standard response only
// when the http.ResponseWriter.Write fails due to unexpected errors.
func writeJSON(w http.ResponseWriter, data GeneralResponse) error {
	response := struct {
		Status  string          `json:"status"`
		Message string          `json:"message"`
		Data    GeneralResponse `json:"data"`
	}{
		Status:  data.Status(),
		Message: data.Message(),
		Data:    data,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(data.Code())
	_, err = w.Write(responseJSON)
	if err != nil {
		return err
	}

	return nil
}
