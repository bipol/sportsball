package handlers

import (
	"fmt"
	"goji.io/pat"
	"net/http"
)

//Hello is just a hello world func
func Hello(w http.ResponseWriter, r *http.Request) {
	name := pat.Param(r, "name")
	fmt.Fprintf(w, "Hello, %s!", name)
}
