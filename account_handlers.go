package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func (s *Server) ListAccounts(w http.ResponseWriter, r *http.Request) {
	rateName := strings.ToUpper(r.URL.Query().Get("rate"))
	accounts, err := ListAccounts(s.db, rateName)

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, accounts, http.StatusOK)
	return
}

func (s *Server) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, _ := strconv.Atoi(vars["id"])
	rateName := strings.ToUpper(r.URL.Query().Get("rate"))
	account, err := GetAccount(s.db, accountID, rateName)

	if err == sql.ErrNoRows {
		respondWithError(w, "not found", http.StatusNotFound)
		return
	}

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, account, http.StatusOK)
	return
}

func (s *Server) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account Account
	json.NewDecoder(r.Body).Decode(&account)

	err := validateRequest(account)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = account.Create(s.db)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, account, http.StatusCreated)
	return
}

func (s *Server) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, _ := strconv.Atoi(vars["id"])
	account, err := GetAccount(s.db, accountID, "")

	err = validateRequest(account)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewDecoder(r.Body).Decode(&account)
	err = account.Update(s.db)

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, account, http.StatusOK)
	return
}

func (s *Server) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, _ := strconv.Atoi(vars["id"])
	account, err := GetAccount(s.db, accountID, "")

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = account.Delete(s.db)

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, nil, http.StatusNoContent)
	return
}
