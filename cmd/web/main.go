package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/crate", snippetCreate)
	log.Print("Starting Server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
