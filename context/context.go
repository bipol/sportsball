package context

import (
	logrus "github.com/Sirupsen/logrus"
	"github.com/bipol/sportsball/config"
	"github.com/bipol/sportsball/storage"
	"github.com/quipo/statsd"
	"time"
)

//AppCtx contains the context of the app to be passed around to handler functions
type AppCtx struct {
	Logger   *logrus.Logger
	Database *storage.DatabaseContext
	Config   *config.Config
	Stats    *statsd.StatsdBuffer
}

//New instantiates a new app context
func New(conf config.Config) (*AppCtx, error) {
	logger := logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)

	database, error := storage.New(conf)
	if error != nil {
		return nil, error
	}

	client := statsd.NewStatsdClient(conf.StatsDHost, conf.StatsDPrefix)
	error = client.CreateSocket()

	if error != nil {
		return nil, error
	}

	collector := statsd.NewStatsdBuffer(time.Second*15, client)

	if error != nil {
		return nil, error
	}

	return &AppCtx{Logger: logger, Database: database, Config: &conf, Stats: collector}, nil
}
