package db

import (
	"database/sql"
	"log"
)

type Command struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Script string `json:"script"`
}

func Database() *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:12345@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Commands (id SERIAL PRIMARY KEY, Name TEXT, Status TEXT, Script TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	return db
}
