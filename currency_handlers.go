package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func (s *Server) ListCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies, err := ListCurrencies(s.db)

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, currencies, http.StatusOK)
	return
}

func (s *Server) GetCurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	currencyName := strings.ToUpper(vars["name"])

	currency, err := GetCurrency(s.db, currencyName)

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, currency, http.StatusOK)
	return
}
