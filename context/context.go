package context

import (
	logrus "github.com/Sirupsen/logrus"
	"github.com/bipol/sportsball/config"
	"github.com/bipol/sportsball/storage"
)

//AppCtx contains the context of the app to be passed around to handler functions
type AppCtx struct {
	Logger   *logrus.Logger
	Database *storage.DatabaseContext
	Config   *conf.Config
}

//New instantiates a new app context
func New(conf config.Config) (*AppCtx, error) {
	logger := logrus.New()

	database, error := storage.New(conf)

	if error != nil {
		return nil, error
	}

	return &AppCtx{logger, database}, nil
}
