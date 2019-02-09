package main

import (
	"testing"

	. "github.com/jonatasbaldin/fin/test"
	"github.com/stretchr/testify/assert"
)

func TestIncomeTransaction(t *testing.T) {
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

	account, _ = GetAccount(s.db, account.ID, "")
	assert.Equal(t, account.Balance, Decimal("100.99"))
}

func TestExpenseTransaction(t *testing.T) {
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
	expense := Transaction{
		Account:     account,
		Description: "My Expense",
		Value:       Decimal("1.00"),
		Type:        "EXPENSE",
		Categories:  []Category{category},
	}
	expense.Create(s.db)

	account, _ = GetAccount(s.db, account.ID, "")
	assert.Equal(t, account.Balance, Decimal("99.99"))
}
