package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	. "github.com/jonatasbaldin/fin/test"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type fakeHttpClient struct {
}

func (client fakeHttpClient) Get(url string) (*http.Response, error) {
	rates := make(map[string]float64)
	rates["BRL"] = 10.99
	rates["ZAR"] = 11.99
	er := ExchangeRates{
		Rates: rates,
		Base:  "USD",
		Date:  "2019-01-28",
	}

	b, err := json.Marshal(er)

	if err != nil {
		log.Fatal("error while marshaling the mocked data", err)
	}

	resp := &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer(b))}

	return resp, nil
}

func TestScrape(t *testing.T) {
	ClearDB(s.db)
	client := fakeHttpClient{}

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

	err := Scrape(s.db, client)
	if err != nil {
		log.Fatal(err)
	}

	curr, err := GetCurrency(s.db, "USD")
	if err != nil {
		log.Fatal(err)
	}

	err = curr.GetLatestRates(s.db)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, len(curr.Rates), 2)
}
