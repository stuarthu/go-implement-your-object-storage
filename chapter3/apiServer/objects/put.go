package objects

import (
	"log"
	"net/http"
)

func put(w http.ResponseWriter, r *http.Request) {
	hash := getHashFromHeader(r)
	if hash == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	e := storeObject(r, hash)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	e = addVersion(r, hash)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
