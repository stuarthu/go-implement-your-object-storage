package objects

import (
	"../heartbeat"
	"io"
	"log"
	"net/http"
)

func put(w http.ResponseWriter, r *http.Request) {
	s := heartbeat.ChooseRandomDataServer()
	if s == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	request, e := http.NewRequest("PUT", "http://"+s+r.URL.Path, r.Body)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := http.Client{}
	nr, e := client.Do(request)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(nr.StatusCode)
	io.Copy(w, nr.Body)
}
