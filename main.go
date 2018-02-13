package main

import (
	"fmt"
	"github.com/bipol/sportsball/config"
	ctx "github.com/bipol/sportsball/context"
	"github.com/bipol/sportsball/handlers"
	"goji.io"
	"goji.io/pat"
	"net/http"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It's alive!")
}

func main() {
	conf := config.New()
	appContext, err := ctx.New(conf)

	if err != nil {
		panic(err)
	}

	mux := goji.NewMux()
	api := handlers.APIMux(appContext)

	//TODO: Auth around API submux
	mux.Handle(pat.New("/api/*"), api)

	mux.HandleFunc(pat.Get("/healthcheck"), healthcheck)

	if err = http.ListenAndServe("localhost:8000", mux); err != nil {
		panic(err)
	}

	appContext.Logger.Info("Server started on 8000")
}
