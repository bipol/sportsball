package storage

import (
	"database/sql"
	"github.com/bipol/sportsball/config"
)

//DatabaseContext contains the connection to the database
type DatabaseContext struct {
	Connection *sql.DB
}

//New instantiates a new database connection
func New(conf config.Config) (*DatabaseContext, error) {
	context := &DatabaseContext{}

	db, err := sql.Open("mysql", conf.DatabaseURL)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	context.Connection = db

	return context, nil
}
