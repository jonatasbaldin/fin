package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID          int             `json:"id"`
	Account     Account         `json:"account"`
	Description string          `json:"description"`
	Value       decimal.Decimal `json:"value"`
	Type        string          `json:"type"`
	Categories  []Category      `json:"categories"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

func (transaction Transaction) MarshalJSON() ([]byte, error) {
	var tmp struct {
		ID          int             `json:"id"`
		Description string          `json:"description"`
		Value       decimal.Decimal `json:"value"`
		Type        string          `json:"type"`
		Categories  []Category      `json:"categories"`
		CreatedAt   string          `json:"created_at"`
		UpdatedAt   string          `json:"updated_at"`
	}

	tmp.ID = transaction.ID
	tmp.Description = transaction.Description
	tmp.Value = transaction.Value
	tmp.Type = transaction.Type
	tmp.Categories = transaction.Categories
	tmp.CreatedAt = transaction.CreatedAt
	tmp.UpdatedAt = transaction.UpdatedAt

	return json.Marshal(&tmp)
}

func ListTransactions(db *sql.DB, accountId int) ([]Transaction, error) {
	rows, err := db.Query("SELECT id, account_id, description, value, type, created_at, updated_at FROM transactions WHERE account_id = $1", accountId)

	if err != nil {
		return nil, err
	}

	transactions := []Transaction{}

	for rows.Next() {
		var transaction Transaction
		errScan := rows.Scan(
			&transaction.ID,
			&transaction.Account.ID,
			&transaction.Description,
			&transaction.Value,
			&transaction.Type,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)

		if errScan != nil {
			return nil, err
		}

		errCat := transaction.GetRelatedCategories(db)

		if errCat != nil {
			return nil, errCat
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (transaction *Transaction) Create(db *sql.DB) (err error) {
	err = transaction.Validate()
	if err != nil {
		return
	}

	createdAt := time.Now()

	tx, err := db.Begin()
	if err != nil {
		return
	}

	err = tx.QueryRow(
		"INSERT INTO transactions(account_id, description, value, type, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at",
		transaction.Account.ID,
		transaction.Description,
		transaction.Value,
		transaction.Type,
		createdAt,
		createdAt,
	).Scan(&transaction.ID, &transaction.CreatedAt, &transaction.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return
	}

	err = transaction.CreateRelatedCategories(db, tx)
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return
}

func GetTransaction(db *sql.DB, id int) (transaction Transaction, err error) {
	err = db.QueryRow(
		"SELECT id, description, value, type, created_at, updated_at FROM transactions WHERE id = $1", id,
	).Scan(
		&transaction.ID,
		&transaction.Description,
		&transaction.Value,
		&transaction.Type,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return
	}

	err = transaction.GetRelatedCategories(db)

	if err != nil {
		return
	}

	return
}

func (transaction *Transaction) Update(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	updatedAt := time.Now()

	_, err = tx.Exec(
		"UPDATE transactions SET description = $1, value = $2, type = $3, updated_at = $4 WHERE id = $5",
		transaction.Description,
		transaction.Value,
		transaction.Type,
		updatedAt,
		transaction.ID,
	)
	if err != nil {
		tx.Rollback()
		return
	}

	err = transaction.DeleteRelatedCategories(db, tx)
	if err != nil {
		tx.Rollback()
		return
	}

	err = transaction.CreateRelatedCategories(db, tx)
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	err = db.QueryRow(
		"SELECT name FROM accounts WHERE id = $1",
		transaction.Account.ID,
	).Scan(&transaction.Account.Name)

	if err != nil {
		return
	}

	return
}

func (transaction *Transaction) Delete(db *sql.DB) (err error) {
	_, err = db.Exec(
		"DELETE FROM transactions WHERE id = $1;",
		transaction.ID,
	)

	if err != nil {
		return
	}

	return
}

func (transaction *Transaction) GetRelatedCategories(db *sql.DB) (err error) {
	rows, err := db.Query(
		"SELECT category_id FROM transactions_categories WHERE transaction_id = $1",
		transaction.ID,
	)

	if err != nil {
		return
	}

	for rows.Next() {
		var categoryId *int
		var category Category

		errScan := rows.Scan(
			&categoryId,
		)
		if errScan != nil {
			return
		}

		category, errCat := GetCategory(db, *categoryId)
		if errCat != nil {
			return
		}

		transaction.Categories = append(transaction.Categories, category)
	}

	return
}

func (transaction *Transaction) CreateRelatedCategories(db *sql.DB, tx *sql.Tx) (err error) {
	for catIndex, structCategory := range transaction.Categories {
		category, categoryErr := GetCategory(db, structCategory.ID)

		if categoryErr != nil {
			tx.Rollback()
			return errors.New(fmt.Sprintf("category %d not found", structCategory.ID))
		}

		transaction.Categories[catIndex] = category

		_, insertErr := tx.Exec(
			"INSERT INTO transactions_categories(transaction_id, category_id) VALUES($1, $2)",
			transaction.ID,
			category.ID,
		)
		if insertErr != nil {
			tx.Rollback()
			return insertErr
		}
	}

	return
}

func (transaction *Transaction) DeleteRelatedCategories(db *sql.DB, tx *sql.Tx) (err error) {
	_, err = tx.Exec(
		"DELETE FROM transactions_categories WHERE transaction_id = $1",
		transaction.ID,
	)

	if err != nil {
		tx.Rollback()
		return
	}

	return
}

func (transaction Transaction) Validate() (err error) {
	typeCheck := regexp.MustCompile(`INCOME|EXPENSE`)
	if !typeCheck.MatchString(transaction.Type) {
		err = errors.New("field 'type' must be 'INCOME' or 'EXPENSE'")
	}

	if transaction.Value.LessThanOrEqual(decimal.NewFromFloat(0)) != false {
		err = errors.New("field 'value' must be more than 0")
	}

	valueCheck := regexp.MustCompile(`^\d*(\.\d{1,2}|\d)$`)
	if !valueCheck.MatchString(transaction.Value.String()) {
		err = errors.New("field 'initial_balance' must be like 1.99")
	}

	if len(transaction.Categories) <= 0 {
		err = errors.New("field 'categories' must not be empty")
	}

	return
}
