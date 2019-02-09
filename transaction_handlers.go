package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) ListTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, _ := strconv.Atoi(vars["account_id"])

	transactions, err := ListTransactions(s.db, accountID)

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, transactions, http.StatusOK)
	return
}

func (s *Server) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, _ := strconv.Atoi(vars["account_id"])

	var transaction Transaction
	json.NewDecoder(r.Body).Decode(&transaction)

	err := validateRequest(transaction)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction.Account.ID = accountID

	err = transaction.Create(s.db)

	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, transaction, http.StatusCreated)
	return
}

func (s *Server) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID, _ := strconv.Atoi(vars["id"])

	transaction, err := GetTransaction(s.db, transactionID)

	if err == sql.ErrNoRows {
		respondWithError(w, "not found", http.StatusNotFound)
		return
	}

	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, transaction, http.StatusOK)
	return
}

func (s *Server) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, _ := strconv.Atoi(vars["account_id"])
	transactionID, _ := strconv.Atoi(vars["id"])
	transaction, err := GetTransaction(s.db, transactionID)
	transaction.Account.ID = accountID

	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewDecoder(r.Body).Decode(&transaction)

	err = validateRequest(transaction)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = transaction.Update(s.db)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, transaction, http.StatusOK)
	return
}

func (s *Server) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID, _ := strconv.Atoi(vars["id"])
	transaction, err := GetTransaction(s.db, transactionID)

	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = transaction.Delete(s.db)

	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, nil, http.StatusNoContent)
	return
}
