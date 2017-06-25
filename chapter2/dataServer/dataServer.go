package main

import (
	"../lib/cfg"
	"./locate"
	"./objects"
	"log"
	"net/http"
)

func objectsHandler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		objects.Put(w, r)
		return
	}
	if m == http.MethodGet {
		objects.Get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func main() {
	go locate.StartLocate()
	http.HandleFunc("/", objectsHandler)
	log.Fatal(http.ListenAndServe(cfg.Config.ListenAddress, nil))
}
