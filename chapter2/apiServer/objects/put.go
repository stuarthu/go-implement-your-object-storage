package objects

import (
	"../heartbeat"
	"io"
	"log"
	"net/http"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	s := heartbeat.ChooseRandomDataServer()
	if s == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	request, e := http.NewRequest("PUT", "http://"+s+"/objects/"+object, r.Body)
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
