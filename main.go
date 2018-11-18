package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	err := http.ListenAndServe(":9876", mux)
	if err != nil {
		log.Fatal(err)
	}
}
