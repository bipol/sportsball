package config

import (
	"fmt"
	"os"
)

//Config contains environment information
type Config struct {
	DatabaseHost string
	DatabasePort string
	DatabaseName string
	DatabaseUser string
	DatabasePass string
	DatabaseURL  string
	StatsDHost   string
	StatsDPrefix string
}

//New instantiates a new config object
func New() Config {

	// Database related
	databaseHost := os.Getenv("DB_HOST")
	databasePort := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_NAME")
	databaseUser := os.Getenv("DB_USER")
	databasePassword := os.Getenv("DB_PASS")

	databaseURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", databaseUser, databasePassword, databaseHost, databasePort, databaseName)

	// Stats
	statsdHost := os.Getenv("STATSD_HOST")
	statsdPrefix := os.Getenv("STATSD_PREFIX")

	conf := Config{
		DatabaseHost: databaseHost,
		DatabasePort: databasePort,
		DatabaseName: databaseName,
		DatabaseUser: databaseUser,
		DatabasePass: databasePassword,
		DatabaseURL:  databaseURL,
		StatsDHost:   statsdHost,
		StatsDPrefix: statsdPrefix,
	}

	return conf
}
