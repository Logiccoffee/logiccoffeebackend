package main

import (
	"net/http"

	"github.com/gocroot/route"
)

func main() {
	http.HandleFunc("/", route.URL)
	http.ListenAndServe(":8080", nil)
}
