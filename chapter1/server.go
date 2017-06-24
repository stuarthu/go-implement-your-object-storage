package main

import (
	. "./handler"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		Put(w, r)
		return
	}
	if m == http.MethodGet {
		Get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":12345", nil)
}
