package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func initializeScrape(db *sql.DB) {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	err := Scrape(db, httpClient)
	if err != nil {
		log.Fatal(err)
	}
}

func isCurrenciesEmpty(db *sql.DB) (ok bool) {
	currencies, err := ListCurrencies(db)
	if err != nil {
		log.Fatal(err)
	}

	if len(currencies) == 0 {
		return true
	}

	return
}

func main() {
	serve := flag.Bool("serve", false, "initialize server")
	scrape := flag.Bool("scrape", false, "initialize scrapper")
	migrate := flag.Bool("migrate", false, "migrate database")
	flag.Parse()

	if len(os.Args) > 1 {
		if flag.NFlag() != 1 {
			fmt.Println("pass just one argument")
			flag.Usage()
			os.Exit(1)
		}

		s := Server{}
		s.Initialize()

		if *serve {
			if isCurrenciesEmpty(s.db) {
				initializeScrape(s.db)
			}
			s.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
		}

		if *migrate {
			s.migrate.Up()
		}

		if *scrape {
			initializeScrape(s.db)
		}

	} else {
		flag.Usage()
	}
}
