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

func TestEmptyListCategories(t *testing.T) {
	ClearDB(s.db)

	response := Request(s.router, "GET", "/categories", nil)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, response.Body.String(), "[]")
}

func TestCreateCategory(t *testing.T) {
	ClearDB(s.db)

	body := []byte(`{"name": "My Category"}`)
	response := Request(s.router, "POST", "/categories", bytes.NewBuffer(body))
	assert.Equal(t, http.StatusCreated, response.Code)

	var respCategory Category
	json.Unmarshal(response.Body.Bytes(), &respCategory)

	assert.Equal(t, respCategory.Name, "My Category")
}

func TestGetCategory(t *testing.T) {
	ClearDB(s.db)

	category := Category{
		Name: "My Category",
	}
	category.Create(s.db)

	response := Request(s.router, "GET", fmt.Sprintf("/categories/%d", category.ID), nil)
	assert.Equal(t, http.StatusOK, response.Code)

	var respCategory Category
	json.Unmarshal(response.Body.Bytes(), &respCategory)

	assert.Equal(t, respCategory.ID, category.ID)
	assert.Equal(t, respCategory.Name, category.Name)
}

func TestUpdateCategory(t *testing.T) {
	ClearDB(s.db)

	category := Category{
		Name: "My Category",
	}
	category.Create(s.db)

	body := []byte(`{"name": "Edited Category"}`)
	response := Request(s.router, "PATCH", fmt.Sprintf("/categories/%d", category.ID), bytes.NewBuffer(body))
	assert.Equal(t, http.StatusOK, response.Code)

	var respCategory Category
	json.Unmarshal(response.Body.Bytes(), &respCategory)

	assert.Equal(t, respCategory.ID, category.ID)
	assert.Equal(t, respCategory.Name, "Edited Category")
}

func TestDeleteCategoryUsedByTransaction(t *testing.T) {
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

	response := Request(s.router, "DELETE", fmt.Sprintf("/categories/%d", category.ID), nil)
	assert.Equal(t, http.StatusInternalServerError, response.Code)

	var err CustomError
	json.Unmarshal(response.Body.Bytes(), &err)

	assert.Equal(
		t,
		err.Error,
		fmt.Sprintf("category '%d' is being used in one or more transaction, please delete them first", category.ID),
	)
}

func TestDeleteCategory(t *testing.T) {
	ClearDB(s.db)

	category := Category{
		Name: "My Category",
	}
	category.Create(s.db)

	response := Request(s.router, "DELETE", fmt.Sprintf("/categories/%d", category.ID), nil)
	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestValidateCategoryName(t *testing.T) {
	ClearDB(s.db)

	body := []byte(`{}`)
	response := Request(s.router, "POST", "/categories", bytes.NewBuffer(body))
	assert.Equal(t, http.StatusBadRequest, response.Code)

	var err CustomError
	json.Unmarshal(response.Body.Bytes(), &err)

	assert.Equal(t, err.Error, "field `name` must not be empty")
}
