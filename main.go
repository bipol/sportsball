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

//Hello is just a hello world func
func hello(w http.ResponseWriter, r *http.Request) {
	name := pat.Param(r, "name")
	fmt.Fprintf(w, "Hello, %s!", name)
}

func main() {
	conf := config.New()
	appContext, err := ctx.New(conf)

	if err != nil {
		panic(err)
	}

	mux := goji.NewMux()
	api := handlers.APIMux(appContext)

	mux.Handle(pat.New("/api/*"), api)

	mux.HandleFunc(pat.Get("/hello/:name"), hello)

	if err = http.ListenAndServe("localhost:8000", mux); err != nil {
		panic(err)
	}

	appContext.Logger.Info("Server started on 8000")
}
