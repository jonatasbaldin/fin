package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Category struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ListCategories(db *sql.DB) ([]Category, error) {
	rows, err := db.Query("SELECT id, name, created_at, updated_at FROM categories")

	if err != nil {
		return nil, err
	}

	categories := []Category{}

	for rows.Next() {
		var category Category
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (category *Category) Create(db *sql.DB) (err error) {
	createdAt := time.Now()

	err = db.QueryRow(
		"INSERT INTO categories(name, created_at, updated_at) VALUES($1, $2, $3) RETURNING id, created_at, updated_at",
		category.Name,
		createdAt,
		createdAt,
	).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return
	}

	return
}

func GetCategory(db *sql.DB, id int) (category Category, err error) {
	err = db.QueryRow(
		"SELECT id, name, created_at, updated_at FROM categories WHERE id = $1", id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return
	}

	return
}

func (category *Category) Update(db *sql.DB) (err error) {
	updatedAt := time.Now()

	err = db.QueryRow(
		"UPDATE categories SET name = $1, updated_at = $2 WHERE id = $3 RETURNING name, updated_at",
		category.Name,
		updatedAt,
		category.ID,
	).Scan(&category.Name, &category.UpdatedAt)

	if err != nil {
		return
	}

	return
}

func (category *Category) Delete(db *sql.DB) (err error) {
	var count int
	err = db.QueryRow(
		"SELECT count(category_id) FROM transactions_categories WHERE category_id = $1",
		category.ID,
	).Scan(&count)

	if err != nil {
		return
	}

	if count == 0 {
		_, err = db.Exec(
			"DELETE FROM categories WHERE id = $1",
			category.ID,
		)

		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("category '%d' is being used in one or more transaction, please delete them first", category.ID)
		return

	}

	return
}

func (category Category) Validate() (err error) {
	if category.Name == "" {
		err = errors.New("field `name` must not be empty")
	}

	return
}
