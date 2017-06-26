package main

import (
	"./objects"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(":12345", nil))
}
