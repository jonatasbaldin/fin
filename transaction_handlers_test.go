package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	. "github.com/jonatasbaldin/fin/test"
	"github.com/stretchr/testify/assert"
)

func TestEmptyListTransactions(t *testing.T) {
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

	response := Request(s.router, "GET", fmt.Sprintf("/accounts/%d/transactions", account.ID), nil)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, response.Body.String(), "[]")
}

func TestListTransactions(t *testing.T) {
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
	category := Category{
		Name: "Category",
	}
	category.Create(s.db)
	transaction := Transaction{
		Account:     account,
		Description: "My Transaction",
		Value:       Decimal("0.99"),
		Type:        "INCOME",
		Categories:  []Category{category},
	}
	transaction.Create(s.db)

	response := Request(s.router, "GET", fmt.Sprintf("/accounts/%d/transactions", account.ID), nil)
	assert.Equal(t, http.StatusOK, response.Code)

	var respTransactions []Transaction
	json.Unmarshal(response.Body.Bytes(), &respTransactions)

	assert.Equal(t, respTransactions[0].Description, transaction.Description)
	assert.Equal(t, respTransactions[0].Value, transaction.Value)
	assert.Equal(t, respTransactions[0].Type, transaction.Type)
	assert.Equal(t, respTransactions[0].Categories[0].Name, transaction.Categories[0].Name)
	assert.Equal(t, len(respTransactions[0].Categories), len(transaction.Categories))
}

func TestCreateTransactionValidCategories(t *testing.T) {
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
	categoryOne := Category{
		Name: "Category One",
	}
	categoryOne.Create(s.db)
	categoryTwo := Category{
		Name: "Category Two",
	}
	categoryTwo.Create(s.db)

	body := []byte(
		fmt.Sprintf(`{"description": "My Transaction", "value": 0.99, "type": "INCOME", "categories": [{"id": %d}, {"id": %d}]}`, categoryOne.ID, categoryTwo.ID),
	)
	response := Request(s.router, "POST", fmt.Sprintf("/accounts/%d/transactions", account.ID), bytes.NewBuffer(body))
	assert.Equal(t, http.StatusCreated, response.Code)

	var respTransaction Transaction
	json.Unmarshal(response.Body.Bytes(), &respTransaction)

	assert.Equal(t, respTransaction.Description, "My Transaction")
	assert.Equal(t, respTransaction.Value, Decimal("0.99"))
	assert.Equal(t, respTransaction.Type, "INCOME")
	assert.Equal(t, respTransaction.Categories[0].Name, "Category One")
	assert.Equal(t, respTransaction.Categories[1].Name, "Category Two")
	assert.Equal(t, len(respTransaction.Categories), 2)
}

func TestCreateTransactionInvalidCategories(t *testing.T) {
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
	categoryOne := Category{
		Name: "Category One",
	}
	categoryOne.Create(s.db)

	body := []byte(
		fmt.Sprintf(`{"description": "My Transaction", "value": 0.99, "type": "INCOME", "categories": [{"id": %d}, {"id": 3213}]}`, categoryOne.ID),
	)
	response := Request(s.router, "POST", fmt.Sprintf("/accounts/%d/transactions", account.ID), bytes.NewBuffer(body))
	assert.Equal(t, http.StatusBadRequest, response.Code)

	var err CustomError
	json.Unmarshal(response.Body.Bytes(), &err)

	assert.Equal(t, err.Error, "category 3213 not found")

	_, errTransaction := GetTransaction(s.db, 1)
	assert.NotNil(t, errTransaction)
}

func TestCreateTransactionEmptyCategories(t *testing.T) {
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

	body := []byte(
		`{"description": "My Transaction", "value": 0.99, "type": "INCOME"}`,
	)
	response := Request(s.router, "POST", fmt.Sprintf("/accounts/%d/transactions", account.ID), bytes.NewBuffer(body))
	assert.Equal(t, http.StatusBadRequest, response.Code)

	var err CustomError
	json.Unmarshal(response.Body.Bytes(), &err)

	assert.Equal(t, err.Error, "field 'categories' must not be empty")

	_, errTransaction := GetTransaction(s.db, 1)
	assert.NotNil(t, errTransaction)
}

func TestGetTransaction(t *testing.T) {
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
	category := Category{
		Name: "Category",
	}
	category.Create(s.db)
	transaction := Transaction{
		Account:     account,
		Description: "My Transaction",
		Value:       Decimal("0.99"),
		Type:        "INCOME",
		Categories:  []Category{category},
	}
	transaction.Create(s.db)

	response := Request(s.router, "GET", fmt.Sprintf("/accounts/%d/transactions/%d", account.ID, transaction.ID), nil)
	assert.Equal(t, http.StatusOK, response.Code)

	var respTransaction Transaction
	json.Unmarshal(response.Body.Bytes(), &respTransaction)

	assert.Equal(t, respTransaction.Description, transaction.Description)
	assert.Equal(t, respTransaction.Value, transaction.Value)
	assert.Equal(t, respTransaction.Type, transaction.Type)
	assert.Equal(t, respTransaction.Categories[0].Name, transaction.Categories[0].Name)
	assert.Equal(t, len(respTransaction.Categories), len(transaction.Categories))
	assert.Equal(t, respTransaction.CreatedAt, transaction.CreatedAt)
}

func TestGetInexistentTransaction(t *testing.T) {
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

	response := Request(s.router, "GET", fmt.Sprintf("/accounts/%d/transactions/12312", account.ID), nil)
	assert.Equal(t, http.StatusNotFound, response.Code)

	var err CustomError
	json.Unmarshal(response.Body.Bytes(), &err)

	assert.Equal(t, err.Error, "not found")
}

func TestUpdateTransactionValidCategories(t *testing.T) {
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
	categoryOne := Category{
		Name: "Category One",
	}
	categoryOne.Create(s.db)
	transaction := Transaction{
		Account:     account,
		Description: "My Transaction",
		Value:       Decimal("0.99"),
		Type:        "INCOME",
		Categories:  []Category{categoryOne},
	}
	transaction.Create(s.db)
	accountCash := Account{
		Currency:       currency,
		Name:           "Cash",
		InitialBalance: Decimal("100.00"),
	}
	accountCash.Create(s.db)
	categoryTwo := Category{
		Name: "Category Two",
	}
	categoryTwo.Create(s.db)

	body := []byte(fmt.Sprintf(`{"account": {"id": %d}, "description": "Edited Transaction", "value": 1.99, "categories": [{"id": %d}]}`, accountCash.ID, categoryTwo.ID))
	response := Request(s.router, "PATCH", fmt.Sprintf("/accounts/%d/transactions/%d", account.ID, transaction.ID), bytes.NewBuffer(body))
	assert.Equal(t, http.StatusOK, response.Code)

	var respTransaction Transaction
	json.Unmarshal(response.Body.Bytes(), &respTransaction)

	assert.Equal(t, respTransaction.Description, "Edited Transaction")
	assert.Equal(t, respTransaction.Value, Decimal("1.99"))
	assert.Equal(t, respTransaction.Type, "INCOME")
	assert.Equal(t, len(respTransaction.Categories), 1)
	assert.Equal(t, respTransaction.CreatedAt, transaction.CreatedAt)
}

func TestUpdateTransactionInvalidCategories(t *testing.T) {
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
	categoryOne := Category{
		Name: "Category One",
	}
	categoryOne.Create(s.db)
	transaction := Transaction{
		Account:     account,
		Description: "My Transaction",
		Value:       Decimal("0.99"),
		Type:        "INCOME",
		Categories:  []Category{categoryOne},
	}
	transaction.Create(s.db)

	body := []byte(`{"categories": [{"id": 2}]}`)
	response := Request(s.router, "PATCH", fmt.Sprintf("/accounts/%d/transactions/%d", account.ID, transaction.ID), bytes.NewBuffer(body))
	assert.Equal(t, http.StatusBadRequest, response.Code)

	var err CustomError
	json.Unmarshal(response.Body.Bytes(), &err)

	assert.Equal(t, err.Error, "category 2 not found")
}

func TestUpdateTransactionEmptyCategories(t *testing.T) {
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
	categoryOne := Category{
		Name: "Category One",
	}
	categoryOne.Create(s.db)
	transaction := Transaction{
		Account:     account,
		Description: "My Transaction",
		Value:       Decimal("0.99"),
		Type:        "INCOME",
		Categories:  []Category{categoryOne},
	}
	transaction.Create(s.db)

	body := []byte(`{"description": "My Edited Transaction"}`)
	response := Request(s.router, "PATCH", fmt.Sprintf("/accounts/%d/transactions/%d", account.ID, transaction.ID), bytes.NewBuffer(body))
	assert.Equal(t, http.StatusOK, response.Code)

	var respTransaction Transaction
	json.Unmarshal(response.Body.Bytes(), &respTransaction)

	assert.Equal(t, respTransaction.Description, "My Edited Transaction")
	assert.Equal(t, len(respTransaction.Categories), 1)
}

func TestUpdateTransactionDifferentCategory(t *testing.T) {
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
	categoryOne := Category{
		Name: "Category One",
	}
	categoryOne.Create(s.db)
	categoryTwo := Category{
		Name: "Category Two",
	}
	categoryTwo.Create(s.db)
	transaction := Transaction{
		Account:     account,
		Description: "My Transaction",
		Value:       Decimal("0.99"),
		Type:        "INCOME",
		Categories:  []Category{categoryOne},
	}
	transaction.Create(s.db)

	body := []byte(fmt.Sprintf(`{"description": "My Edited Transaction", "categories": [{"id": %d}]}`, categoryTwo.ID))
	response := Request(s.router, "PATCH", fmt.Sprintf("/accounts/%d/transactions/%d", account.ID, transaction.ID), bytes.NewBuffer(body))
	assert.Equal(t, http.StatusOK, response.Code)

	var respTransaction Transaction
	json.Unmarshal(response.Body.Bytes(), &respTransaction)

	assert.Equal(t, len(respTransaction.Categories), 1)
	assert.Equal(t, respTransaction.Categories[0].Name, "Category Two")
}

func TestDeleteTransaction(t *testing.T) {
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
	category := Category{
		Name: "Category",
	}
	category.Create(s.db)
	transaction := Transaction{
		Account:     account,
		Description: "My Transaction",
		Value:       Decimal("0.99"),
		Type:        "EXPENSE",
		Categories:  []Category{category},
	}
	transaction.Create(s.db)

	response := Request(s.router, "DELETE", fmt.Sprintf("/accounts/%d/transactions/%d", account.ID, transaction.ID), nil)
	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestValidateTransactionType(t *testing.T) {
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
	categoryOne := Category{
		Name: "Category One",
	}
	categoryOne.Create(s.db)

	body := []byte(`{"description": "My First Transaction", "value": 0.99, "type": "WHAT", "categories": [{"id": 1}]}`)
	response := Request(s.router, "POST", fmt.Sprintf("/accounts/%d/transactions", account.ID), bytes.NewBuffer(body))
	assert.Equal(t, http.StatusBadRequest, response.Code)

	var err CustomError
	json.Unmarshal(response.Body.Bytes(), &err)

	assert.Equal(t, err.Error, "field 'type' must be 'INCOME' or 'EXPENSE'")
}

func TestValidateTransactionValue(t *testing.T) {
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
	categoryOne := Category{
		Name: "Category One",
	}
	categoryOne.Create(s.db)

	body := []byte(`{"description": "My First Transaction", "type": "EXPENSE", "categories": [{"id": 1}]}`)
	response := Request(s.router, "POST", fmt.Sprintf("/accounts/%d/transactions", account.ID), bytes.NewBuffer(body))
	assert.Equal(t, http.StatusBadRequest, response.Code)

	var err CustomError
	json.Unmarshal(response.Body.Bytes(), &err)

	assert.Equal(t, err.Error, "field 'value' must be more than 0")
}
