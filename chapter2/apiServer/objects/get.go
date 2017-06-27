package objects

import (
	"../locate"
	"io"
	"log"
	"net/http"
	"strings"
)

func Get(w http.ResponseWriter, r *http.Request) {
	s := locate.Locate(strings.Split(r.URL.Path, "/")[2])
	if s == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	request, e := http.NewRequest("GET", "http://"+s+r.URL.Path, r.Body)
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
