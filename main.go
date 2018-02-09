package main

import (
	"fmt"
	logrus "github.com/Sirupsen/logrus"
	"github.com/bipol/sportsball/storage"
	"goji.io"
	"goji.io/pat"
	"net/http"
)

type AppCtx struct {
	Logger   *logrus.Logger
	Database *storage.DatabaseContext
}

func hello(w http.ResponseWriter, r *http.Request) {
	name := pat.Param(r, "name")
	fmt.Fprintf(w, "Hello, %s!", name)
}

func generateAppCtx(filename string) (*AppCtx, error) {
	database, error := storage.New(filename)

	if error != nil {
		return nil, error
	}

	logger := logrus.New()

	return &AppCtx{logger, database}, nil
}

func main() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/hello/:name"), hello)

	http.ListenAndServe("localhost:8000", mux)
}
