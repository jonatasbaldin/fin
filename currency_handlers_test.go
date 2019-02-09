package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	. "github.com/jonatasbaldin/fin/test"
)

var s Server

func TestMain(main *testing.M) {
	dbStr := os.Getenv("DB_TEST")
	s.initializeDB(dbStr)
	s.initializeRoutes()
	EnsureTablesExists(s.db)
	code := main.Run()
	ClearDB(s.db)
	os.Exit(code)
}

func TestEmptyListCurrencies(t *testing.T) {
	ClearDB(s.db)
	response := Request(s.router, "GET", "/currencies", nil)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, response.Body.String(), "[]")
}

func TestListCurrencies(t *testing.T) {
	ClearDB(s.db)
	value, _ := decimal.NewFromString("3.80")
	rate := Rate{
		Name:   "BRL",
		Symbol: "R$",
		Value:  value,
	}
	currency := Currency{
		Name:   "USD",
		Symbol: "$",
		Rates:  []Rate{rate},
	}
	currency.Create(s.db)

	response := Request(s.router, "GET", "/currencies", nil)
	assert.Equal(t, http.StatusOK, response.Code)

	var respCurrencies []Currency
	json.Unmarshal(response.Body.Bytes(), &respCurrencies)

	assert.Equal(t, respCurrencies[0].Name, currency.Name)
	assert.Equal(t, respCurrencies[0].Symbol, currency.Symbol)
	assert.Equal(t, respCurrencies[0].Rates[0].Name, rate.Name)
	assert.Equal(t, respCurrencies[0].Rates[0].Symbol, rate.Symbol)
	assert.Equal(t, respCurrencies[0].Rates[0].Value, rate.Value)
}

func TestGetCurrency(t *testing.T) {
	ClearDB(s.db)
	value, _ := decimal.NewFromString("3.80")
	rate := Rate{
		Name:   "BRL",
		Symbol: "R$",
		Value:  value,
	}
	currency := Currency{
		Name:   "USD",
		Symbol: "$",
		Rates:  []Rate{rate},
	}
	currency.Create(s.db)

	response := Request(s.router, "GET", fmt.Sprintf("/currencies/%s", currency.Name), nil)
	assert.Equal(t, http.StatusOK, response.Code)

	var respCurrency Currency
	json.Unmarshal(response.Body.Bytes(), &respCurrency)

	assert.Equal(t, respCurrency.Name, currency.Name)
	assert.Equal(t, respCurrency.Symbol, currency.Symbol)
	assert.Equal(t, respCurrency.Rates[0].Name, rate.Name)
	assert.Equal(t, respCurrency.Rates[0].Symbol, rate.Symbol)
	assert.Equal(t, respCurrency.Rates[0].Value, rate.Value)
}
