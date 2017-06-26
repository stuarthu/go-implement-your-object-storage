package main

import (
	"../lib/cfg"
	"./locate"
	"./objects"
	"log"
	"net/http"
)

func main() {
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(cfg.Config.ListenAddress, nil))
}
