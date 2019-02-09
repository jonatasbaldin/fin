package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/jonatasbaldin/fin/test"
)

func TestEmptyListAccounts(t *testing.T) {
	ClearDB(s.db)
	response := Request(s.router, "GET", "/accounts", nil)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, response.Body.String(), "[]")
}

func TestListAccounts(t *testing.T) {
	ClearDB(s.db)

	currency := Currency{
		Name: "USD",
	}
	currency.Create(s.db)
	account := Account{
		Currency:       currency,
		Name:           "My Wallet",
		InitialBalance: Decimal("100.00"),
	}
	account.Create(s.db)

	response := Request(s.router, "GET", "/accounts", nil)
	assert.Equal(t, http.StatusOK, response.Code)

	var respAccounts []Account
	json.Unmarshal(response.Body.Bytes(), &respAccounts)

	assert.Equal(t, respAccounts[0].Name, account.Name)
	assert.Equal(t, respAccounts[0].InitialBalance, account.InitialBalance)
	assert.Equal(t, respAccounts[0].CreatedAt, account.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	ClearDB(s.db)

	currency := Currency{
		Name: "USD",
	}
	currency.Create(s.db)
	account := Account{
		Currency:       currency,
		Name:           "My Wallet",
		InitialBalance: Decimal("100.00"),
	}
	account.Create(s.db)

	response := Request(s.router, "GET", fmt.Sprintf("/accounts/%d", account.ID), nil)
	assert.Equal(t, http.StatusOK, response.Code)

	var respAccount Account
	json.Unmarshal(response.Body.Bytes(), &respAccount)

	assert.Equal(t, respAccount.Name, account.Name)
	assert.Equal(t, respAccount.InitialBalance, account.InitialBalance)
	assert.Equal(t, respAccount.CreatedAt, account.CreatedAt)
}

func TestCreateAccount(t *testing.T) {
	ClearDB(s.db)

	currency := Currency{
		Name: "USD",
	}
	currency.Create(s.db)

	body := []byte(
		`{"currency": {"name": "USD"}, "name": "My Wallet", "initial_balance": 100.0}`,
	)
	response := Request(s.router, "POST", "/accounts", bytes.NewBuffer(body))

	assert.Equal(t, http.StatusCreated, response.Code)

	var respAccount Account
	json.Unmarshal(response.Body.Bytes(), &respAccount)

	assert.Equal(t, respAccount.Name, "My Wallet")
	assert.Equal(t, respAccount.InitialBalance, Decimal("100.00"))
	assert.Equal(t, respAccount.Balance, Decimal("100.00"))
	assert.Equal(t, respAccount.Currency.Name, "USD")
}

func TestUpdateAccount(t *testing.T) {
	ClearDB(s.db)

	currency := Currency{
		Name: "USD",
	}
	currency.Create(s.db)
	currencyBrl := Currency{
		Name: "BRL",
	}
	currencyBrl.Create(s.db)
	account := Account{
		Currency:       currency,
		Name:           "My Wallet",
		InitialBalance: Decimal("100.00"),
	}
	account.Create(s.db)

	body := []byte(
		`{"currency": {"name": "BRL"}, "name": "Edited Wallet"}`,
	)
	response := Request(s.router, "PATCH", fmt.Sprintf("/accounts/%d", account.ID), bytes.NewBuffer(body))
	assert.Equal(t, http.StatusOK, response.Code)

	var respAccount Account
	json.Unmarshal(response.Body.Bytes(), &respAccount)

	assert.Equal(t, respAccount.Name, "Edited Wallet")
	assert.Equal(t, respAccount.Currency.Name, "BRL")
}

func TestDeleteAccount(t *testing.T) {
	ClearDB(s.db)

	currency := Currency{
		Name: "USD",
	}
	currency.Create(s.db)
	account := Account{
		Currency:       currency,
		Name:           "My Wallet",
		InitialBalance: Decimal("100.00"),
	}
	account.Create(s.db)

	response := Request(s.router, "DELETE", fmt.Sprintf("/accounts/%d", account.ID), nil)
	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestValidateAccountCurrencyID(t *testing.T) {
	ClearDB(s.db)

	currency := Currency{
		Name: "USD",
	}
	currency.Create(s.db)

	body := []byte(`{"name": "My Wallet", "initial_balance": 100.0}`)
	response := Request(s.router, "POST", "/accounts", bytes.NewBuffer(body))
	assert.Equal(t, http.StatusBadRequest, response.Code)

	var err CustomError
	json.Unmarshal(response.Body.Bytes(), &err)

	assert.Equal(t, err.Error, "field 'currency.name' must not be empty")
}

func TestValidateAccountName(t *testing.T) {
	ClearDB(s.db)

	currency := Currency{
		Name: "USD",
	}
	currency.Create(s.db)

	body := []byte(`{"currency": {"name": "USD"}, "initial_balance": 100.0}`)
	response := Request(s.router, "POST", "/accounts", bytes.NewBuffer(body))
	assert.Equal(t, http.StatusBadRequest, response.Code)

	var err CustomError
	json.Unmarshal(response.Body.Bytes(), &err)

	assert.Equal(t, err.Error, "field 'name' must not be empty")
}

func TestValidateAccountInitialBalance(t *testing.T) {
	ClearDB(s.db)

	currency := Currency{
		Name: "USD",
	}
	currency.Create(s.db)

	body := []byte(`{"currency": {"name": "USD"}, "name": "My Wallet"}`)
	response := Request(s.router, "POST", "/accounts", bytes.NewBuffer(body))
	assert.Equal(t, http.StatusBadRequest, response.Code)

	var err CustomError
	json.Unmarshal(response.Body.Bytes(), &err)

	assert.Equal(t, err.Error, "field 'initial_balance' must be more than 0")
}
