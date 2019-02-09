package main

import (
	"encoding/json"
	"net/http"
)

type RequestValidation interface {
	Validate() error
}

func validateRequest(v RequestValidation) (err error) {
	err = v.Validate()
	return
}

func respondWithError(w http.ResponseWriter, errMessage string, statusCode int) {
	message := map[string]string{"error": errMessage}
	response, _ := json.Marshal(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func respondWithJSON(w http.ResponseWriter, payload interface{}, statusCode int) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}
