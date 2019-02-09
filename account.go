package main

import (
	"database/sql"
	"errors"
	"regexp"
	"time"

	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID             int             `json:"id"`
	Currency       Currency        `json:"currency"`
	Name           string          `json:"name"`
	InitialBalance decimal.Decimal `json:"initial_balance"`
	Balance        decimal.Decimal `json:"balance"`
	CreatedAt      string          `json:"created_at"`
	UpdatedAt      string          `json:"updated_at"`
}

func (account *Account) getBalance(db *sql.DB, rateName string) error {
	var income decimal.Decimal
	var expense decimal.Decimal

	err := db.QueryRow(
		"SELECT value FROM transactions WHERE account_id = $1 AND type = 'INCOME'", account.ID,
	).Scan(&income)

	if err == sql.ErrNoRows {
		account.Balance = account.InitialBalance
	} else {
		account.Balance = decimal.Sum(account.InitialBalance, income)
	}

	err = db.QueryRow(
		"SELECT value FROM transactions WHERE account_id = $1 AND type = 'EXPENSE'", account.ID,
	).Scan(&expense)

	if err == sql.ErrNoRows {
		account.Balance = account.Balance
	} else {
		account.Balance = account.Balance.Sub(expense)
	}

	if rateName != "" && account.Currency.Name != rateName {
		errCalc := account.calculateInitialBalance(db, rateName)
		if errCalc != nil {
			return errCalc
		}

		errCalc = account.calculateBalance(db, rateName)
		if errCalc != nil {
			return errCalc
		}
	}

	return nil
}

func (account *Account) calculateInitialBalance(db *sql.DB, rateName string) (err error) {
	rate, err := account.Currency.GetRate(db, rateName)
	account.InitialBalance = account.InitialBalance.Mul(rate.Value).Truncate(2)
	return
}

func (account *Account) calculateBalance(db *sql.DB, rateName string) (err error) {
	rate, err := account.Currency.GetRate(db, rateName)
	account.Balance = account.Balance.Mul(rate.Value).Truncate(2)
	return
}

func ListAccounts(db *sql.DB, rateName string) ([]Account, error) {
	rows, err := db.Query(
		`SELECT a.id, c.name, a.name, a.initial_balance, a.created_at, a.updated_at
	     FROM accounts a
		 INNER JOIN currencies c ON (a.currency_name = c.name)`,
	)

	if err != nil {
		return nil, err
	}

	accounts := []Account{}

	for rows.Next() {
		var account Account
		errScan := rows.Scan(
			&account.ID,
			&account.Currency.Name,
			&account.Name,
			&account.InitialBalance,
			&account.CreatedAt,
			&account.UpdatedAt,
		)

		if errScan != nil {
			return nil, errScan
		}

		account.Currency, errScan = GetCurrency(db, account.Currency.Name)

		if errScan != nil {
			return nil, errScan
		}

		errBalance := account.getBalance(db, rateName)

		if errBalance != nil {
			return nil, errBalance
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func GetAccount(db *sql.DB, id int, rateName string) (account Account, err error) {
	row := db.QueryRow(
		`SELECT a.id, c.name, a.name, initial_balance, a.created_at, a.updated_at
		FROM accounts a INNER JOIN currencies c ON (a.currency_name = c.name)
		WHERE a.id = $1`, id,
	)

	err = row.Scan(
		&account.ID,
		&account.Currency.Name,
		&account.Name,
		&account.InitialBalance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		return
	}

	account.Currency, err = GetCurrency(db, account.Currency.Name)
	if err != nil {
		return
	}

	err = account.getBalance(db, rateName)

	if err != nil {
		return
	}

	return
}

func (account *Account) Create(db *sql.DB) (err error) {
	createdAt := time.Now()

	account.Currency, err = GetCurrency(db, account.Currency.Name)
	if err != nil {
		return
	}

	err = db.QueryRow(
		`INSERT INTO accounts(currency_name, name, initial_balance, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at, updated_at`,
		account.Currency.Name,
		account.Name,
		account.InitialBalance,
		createdAt,
		createdAt,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		return
	}

	err = account.getBalance(db, account.Currency.Name)
	if err != nil {
		return
	}

	return
}

func (account *Account) Update(db *sql.DB) (err error) {
	updatedAt := time.Now()

	err = db.QueryRow(
		`UPDATE accounts
		SET name = $1, updated_at = $2
		WHERE id = $3
		RETURNING updated_at`,
		account.Name,
		updatedAt,
		account.ID,
	).Scan(&account.UpdatedAt)

	if err != nil {
		return
	}

	return
}

func (account *Account) Delete(db *sql.DB) (err error) {
	_, err = db.Exec(
		"DELETE FROM accounts WHERE id = $1;",
		account.ID,
	)

	if err != nil {
		return
	}

	return
}

func (account Account) Validate() (err error) {
	if account.Currency.Name == "" {
		err = errors.New("field 'currency.name' must not be empty")
	}

	if account.Name == "" {
		err = errors.New("field 'name' must not be empty")
	}

	if account.InitialBalance.LessThanOrEqual(decimal.NewFromFloat(0)) != false {
		err = errors.New("field 'initial_balance' must be more than 0")
	}

	valueCheck := regexp.MustCompile(`^\d*(\.\d{1,2}|\d)$`)
	if !valueCheck.MatchString(account.InitialBalance.String()) {
		err = errors.New("field 'initial_balance' must be like 1.99")
	}

	return
}
