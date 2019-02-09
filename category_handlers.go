package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := ListCategories(s.db)

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, categories, http.StatusOK)
	return
}

func (s *Server) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	json.NewDecoder(r.Body).Decode(&category)

	err := validateRequest(category)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = category.Create(s.db)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, category, http.StatusCreated)
	return
}

func (s *Server) GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID, _ := strconv.Atoi(vars["id"])
	category, err := GetCategory(s.db, categoryID)

	if err == sql.ErrNoRows {
		respondWithError(w, "not found", http.StatusNotFound)
		return
	}

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, category, http.StatusOK)
	return
}

func (s *Server) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID, _ := strconv.Atoi(vars["id"])
	category, err := GetCategory(s.db, categoryID)

	err = validateRequest(category)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewDecoder(r.Body).Decode(&category)
	err = category.Update(s.db)

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, category, http.StatusOK)
	return
}

func (s *Server) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID, _ := strconv.Atoi(vars["id"])
	category, err := GetCategory(s.db, categoryID)

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = category.Delete(s.db)

	if err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, nil, http.StatusNoContent)
	return
}
