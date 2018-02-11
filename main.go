package main

import (
	"github.com/bipol/sportsball/config"
	ctx "github.com/bipol/sportsball/context"
	"github.com/bipol/sportsball/handlers"
	"goji.io"
	"goji.io/pat"
	"net/http"
)

func main() {
	conf := config.New()
	appContext, err := ctx.New(conf)

	if err != nil {
		panic(err)
	}

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/hello/:name"), handlers.Hello)

	http.ListenAndServe("localhost:8000", mux)
}
