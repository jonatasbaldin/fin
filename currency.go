package main

import (
	"database/sql"
	"errors"
	"regexp"
	"time"

	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

var supportedRatesAndSymbols = map[string]string{
	"EUR": "€",
	"NZD": "$",
	"CAD": "$",
	"MXN": "$",
	"AUD": "﷼",
	"CNY": "¥",
	"PHP": "₱",
	"GBP": "£",
	"CZK": "Kč",
	"USD": "$",
	"SEK": "kr",
	"NOK": "kr",
	"TRY": "₺",
	"IDR": "Rp",
	"ZAR": "R",
	"MYR": "RM",
	"HKD": "$",
	"HUF": "Ft",
	"ISK": "kr",
	"HRK": "kn",
	"JPY": "¥",
	"BGN": "лв",
	"SGD": "$",
	"RUB": "₽",
	"RON": "lei",
	"CHF": "CHF",
	"DKK": "kr",
	"INR": "₹",
	"KRW": "₩",
	"THB": "฿",
	"BRL": "R$",
	"PLN": "zł",
	"ILS": "₪",
}

type Currency struct {
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	Rates     []Rate `json:"rates"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Rate struct {
	ID        int             `json:"-"`
	Name      string          `json:"name"`
	Symbol    string          `json:"symbol"`
	Value     decimal.Decimal `json:"value"`
	CreatedAt string          `json:"-"`
	UpdatedAt string          `json:"-"`
}

func (currency *Currency) Create(db *sql.DB) (err error) {
	createdAt := time.Now()

	tx, err := db.Begin()
	if err != nil {
		return
	}

	_, err = tx.Exec(
		"INSERT INTO currencies(name, symbol, created_at, updated_at) VALUES ($1, $2, $3, $4)",
		currency.Name,
		currency.Symbol,
		createdAt,
		createdAt,
	)

	if err != nil {
		tx.Rollback()
		return
	}

	err = currency.CreateRates(db, tx)
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

func (currency *Currency) CreateRates(db *sql.DB, tx *sql.Tx) (err error) {
	createdAt := time.Now()

	// always add new rates, because for we always get the latest one with GetLatestRates
	for _, rate := range currency.Rates {
		rateErr := rate.Validate()
		if rateErr != nil {
			return rateErr
		}

		rateErr = tx.QueryRow(
			"INSERT INTO rates(currency_name, name, symbol, value, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
			currency.Name,
			rate.Name,
			rate.Symbol,
			rate.Value,
			createdAt,
			createdAt,
		).Scan(&rate.ID)

		if rateErr != nil {
			return rateErr
		}
	}

	return
}

func ListCurrencies(db *sql.DB) ([]Currency, error) {
	rows, err := db.Query("SELECT name, symbol, created_at, updated_at FROM currencies")

	if err != nil {
		return nil, err
	}

	currencies := []Currency{}

	for rows.Next() {
		var currency Currency
		errScan := rows.Scan(
			&currency.Name,
			&currency.Symbol,
			&currency.CreatedAt,
			&currency.UpdatedAt,
		)

		if errScan != nil {
			return nil, err
		}

		errRates := currency.GetLatestRates(db)
		if errRates != nil {
			return nil, err
		}

		currencies = append(currencies, currency)
	}

	return currencies, nil
}

func (currency *Currency) GetLatestRates(db *sql.DB) (err error) {
	currency.CleanRates()

	// Gets the latest Rate for each currency
	rows, err := db.Query(
		`SELECT r1.id, r1.name , r1.symbol, r1.value, r1.created_at, r1.updated_at
		  FROM (
			SELECT name, currency_name, MAX(created_at) AS created_at
			FROM rates
			GROUP BY name, currency_name
		  ) r2
		  JOIN rates r1
		  ON r1.created_at = r2.created_at AND r1.currency_name = r2.currency_name AND r1.name = r2.name
		  WHERE r1.currency_name = $1`,
		currency.Name,
	)

	if err != nil {
		return
	}

	for rows.Next() {
		var rate Rate

		errScan := rows.Scan(
			&rate.ID,
			&rate.Name,
			&rate.Symbol,
			&rate.Value,
			&rate.CreatedAt,
			&rate.UpdatedAt,
		)
		if errScan != nil {
			return
		}

		currency.Rates = append(currency.Rates, rate)
	}

	return
}

func (currency *Currency) GetRate(db *sql.DB, rateName string) (rate Rate, err error) {
	for _, r := range currency.Rates {
		if r.Name == rateName {
			return r, nil
		}
	}

	return
}

func GetCurrency(db *sql.DB, name string) (currency Currency, err error) {
	err = db.QueryRow(
		"SELECT name, symbol, created_at, updated_at FROM currencies WHERE name = $1", name,
	).Scan(
		&currency.Name,
		&currency.Symbol,
		&currency.CreatedAt,
		&currency.UpdatedAt,
	)

	if err != nil {
		return
	}

	err = currency.GetLatestRates(db)
	if err != nil {
		return
	}

	return
}

// CleanRates is used to clean the currency.Rates
// used when adding new Rates to a Currency, to avoid all Rates getting their time updated
// and when getting the latest rates, to avoid having more than one
func (currency *Currency) CleanRates() {
	currency.Rates = []Rate{}
}

func (rate Rate) Validate() (err error) {
	valueCheck := regexp.MustCompile(`^\d*(\.\d{1,2}|\d)$`)
	if !valueCheck.MatchString(rate.Value.String()) {
		err = errors.New("field 'value' must be like 1.99")
	}

	return
}
