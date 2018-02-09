package storage

import (
	"database/sql"
)

type DatabaseContext struct {
	Connection *sql.DB
}

func New(filename string) (*DatabaseContext, error) {
	context := &DatabaseContext{}

	db, err := sql.Open("sqlite3", filename)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	context.Connection = db

	return context, nil
}
