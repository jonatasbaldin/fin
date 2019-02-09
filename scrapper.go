package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

const exchangeRateAPIURL string = "https://api.exchangeratesapi.io"

type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

type httpClient interface {
	Get(url string) (resp *http.Response, err error)
}

func Scrape(db *sql.DB, httpClient httpClient) (err error) {
	log.Print("Starting scrapper")

	for currencyName, currencySymbol := range supportedRatesAndSymbols {
		currency, err := GetCurrency(db, currencyName)
		if err != nil {
			currency = Currency{
				Name:   currencyName,
				Symbol: currencySymbol,
			}
			err = currency.Create(db)
			if err != nil {
				return err
			}
		}

		var er ExchangeRates
		er.getExchangeRates(currency.Name, httpClient)

		currency.CleanRates()

		for key, value := range er.Rates {
			rate := Rate{
				Name:   key,
				Symbol: supportedRatesAndSymbols[key],
				Value:  decimal.NewFromFloat(value).Truncate(2),
			}
			currency.Rates = append(currency.Rates, rate)
		}

		err = createRates(db, &currency)
		if err != nil {
			return err
		}

	}

	log.Print("Finishing scrapper")

	return
}

func (er *ExchangeRates) getExchangeRates(currencyName string, httpClient httpClient) {
	url := buildURL(currencyName)

	resp, err := httpClient.Get(url)
	if err != nil {
		// TODO: should I panic here?
		panic(err.Error())
	}
	defer resp.Body.Close()

	log.Print(fmt.Sprintf("Getting ExchageRate from %s", url))

	json.NewDecoder(resp.Body).Decode(er)
}

func buildURL(currencyName string) string {
	return exchangeRateAPIURL + fmt.Sprintf("/latest?base=%s", currencyName)
}

func createRates(db *sql.DB, currency *Currency) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = currency.CreateRates(db, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return
}
