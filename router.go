package main

import "github.com/gorilla/mux"

func (s *Server) initializeRoutes() {
	s.router = mux.NewRouter()

	s.router.HandleFunc("/accounts", s.ListAccounts).Methods("GET")
	s.router.HandleFunc("/accounts", s.CreateAccount).Methods("POST")
	s.router.HandleFunc("/accounts/{id:[0-9]+}", s.GetAccount).Methods("GET")
	s.router.HandleFunc("/accounts/{id:[0-9]+}", s.UpdateAccount).Methods("PATCH")
	s.router.HandleFunc("/accounts/{id:[0-9]+}", s.DeleteAccount).Methods("DELETE")
	s.router.HandleFunc("/accounts/{account_id:[0-9]+}/transactions", s.ListTransactions).Methods("GET")
	s.router.HandleFunc("/accounts/{account_id:[0-9]+}/transactions", s.CreateTransaction).Methods("POST")
	s.router.HandleFunc("/accounts/{account_id:[0-9]+}/transactions/{id:[0-9]+}", s.GetTransaction).Methods("GET")
	s.router.HandleFunc("/accounts/{account_id:[0-9]+}/transactions/{id:[0-9]+}", s.UpdateTransaction).Methods("PATCH")
	s.router.HandleFunc("/accounts/{account_id:[0-9]+}/transactions/{id:[0-9]+}", s.DeleteTransaction).Methods("DELETE")
	s.router.HandleFunc("/categories", s.ListCategories).Methods("GET")
	s.router.HandleFunc("/categories", s.CreateCategory).Methods("POST")
	s.router.HandleFunc("/categories/{id:[0-9]+}", s.GetCategory).Methods("GET")
	s.router.HandleFunc("/categories/{id:[0-9]+}", s.UpdateCategory).Methods("PATCH")
	s.router.HandleFunc("/categories/{id:[0-9]+}", s.DeleteCategory).Methods("DELETE")
	s.router.HandleFunc("/currencies", s.ListCurrencies).Methods("GET")
	s.router.HandleFunc("/currencies/{name:[a-zA-Z]{3}}", s.GetCurrency).Methods("GET")
}
