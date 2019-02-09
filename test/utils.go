package test

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

const currenciesTableCreation = `CREATE TABLE IF NOT EXISTS currencies
(
name varchar(255) primary key,
symbol varchar(255) not null,
created_at timestamp not null,
updated_at timestamp not null
)`

const ratesTableCreation = `CREATE TABLE IF NOT EXISTS rates
(
id serial primary key,
currency_name varchar(255) references currencies (name),
name varchar(255) not null,
symbol varchar(255) not null,
value numeric(12,2) not null,
created_at timestamp not null,
updated_at timestamp not null
)
`
const accountsTableCreation = `CREATE TABLE IF NOT EXISTS accounts
(
id serial primary key,
currency_name varchar(255) references currencies (name),
name varchar(255) not null,
initial_balance real not null,
created_at timestamp not null,
updated_at timestamp not null
)`

const transactionsTableCreation = `CREATE TABLE IF NOT EXISTS transactions
(
id serial primary key,
account_id int references accounts (id),
description varchar(255),
value numeric(12,2) not null,
type varchar(255) not null,
created_at timestamp not null,
updated_at timestamp not null
)`

const categoriesTableCreation = `CREATE TABLE IF NOT EXISTS categories
(
id serial primary key,
name varchar(255),
created_at timestamp not null,
updated_at timestamp not null
)`

const transactionsCategoriesTableCreation = ` CREATE TABLE IF NOT EXISTS transactions_categories 
(
transaction_id int REFERENCES transactions ON DELETE CASCADE,
category_id int REFERENCES categories,
PRIMARY KEY (transaction_id, category_id)
)
`

func EnsureTablesExists(db *sql.DB) {
	var err error

	if _, err = db.Exec(currenciesTableCreation); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(accountsTableCreation); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(transactionsTableCreation); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(categoriesTableCreation); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(categoriesTableCreation); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(transactionsCategoriesTableCreation); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(ratesTableCreation); err != nil {
		log.Fatal(err)
	}
}

func ClearDB(db *sql.DB) {
	db.Exec("DELETE FROM transactions_categories")

	db.Exec("DELETE FROM categories")
	// db.Exec("ALTER SEQUENCE categories_id_seq RESTART WITH 1")

	db.Exec("DELETE FROM transactions")
	// db.Exec("ALTER SEQUENCE transactions_id_seq RESTART WITH 1")

	db.Exec("DELETE FROM accounts")
	// db.Exec("ALTER SEQUENCE accounts_id_seq RESTART WITH 1")

	db.Exec("DELETE FROM rates")
	// db.Exec("ALTER SEQUENCE rates_id_seq RESTART WITH 1")

	db.Exec("DELETE FROM currencies")
	// db.Exec("ALTER SEQUENCE currencies_id_seq RESTART WITH 1")
}

func Request(router *mux.Router, method string, path string, body io.Reader) (responseRecorder *httptest.ResponseRecorder) {
	req, err := http.NewRequest(method, path, body)

	if err != nil {
		log.Fatal("Could create the HTTP request")
	}

	responseRecorder = httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)

	return
}

func Decimal(value string) decimal.Decimal {
	dec, _ := decimal.NewFromString(value)

	return dec
}

type CustomError struct {
	Error string
}
