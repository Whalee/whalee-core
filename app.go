package main

import (
	"log"
	"net/http"
)

func main() {
	ReadConfig()
	router := NewRouter()
	log.Println("Starting server on http://localhost:4200")
	log.Fatal(http.ListenAndServe(":4200", router))
}
