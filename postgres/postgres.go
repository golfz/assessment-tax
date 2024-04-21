package postgres

import (
	"database/sql"
	"log"
)

func New(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Printf("unable to open database connection: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Printf("unable to connect to database: %v", err)
		return nil, err
	}

	return db, nil
}
