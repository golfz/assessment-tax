package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Postgres struct {
	Db *sql.DB
}

func New(databaseURL string) (*Postgres, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Printf("unable to open database connection: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Printf("unable to connect to database: %v", err)
		return nil, err
	}

	return &Postgres{Db: db}, nil
}
