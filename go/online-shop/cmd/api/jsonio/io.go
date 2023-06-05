package jsonio

import (
	"encoding/json"
	"net/http"
)

// Write is a utility function to write arbitrary data to the JSON response.
// The response struct always obeys a standard struct providing essential status
// and message to the client. It may fail to write a standard response only when
// the http.ResponseWriter.Write fails due to unexpected errors.
func Write(w http.ResponseWriter, data Response) error {
	responseJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(data.Code)
	_, err = w.Write(responseJSON)
	if err != nil {
		return err
	}

	return nil
}
