package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Server struct {
	db      *sql.DB
	router  *mux.Router
	migrate *migrate.Migrate
}

func (s *Server) initializeDB(dbStr string) {
	var err error
	s.db, err = sql.Open("postgres", dbStr)

	if err != nil {
		log.Fatal(err)
	}

	if err = s.db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) initializeMigrate() {
	driver, err := postgres.WithInstance(s.db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	s.migrate, err = migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) migrateUp() {
	err := s.migrate.Up()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) Initialize() {
	dbStr := os.Getenv("DB")
	s.initializeDB(dbStr)
	s.initializeRoutes()
	s.initializeMigrate()
}

func (s *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, s.router))
}
