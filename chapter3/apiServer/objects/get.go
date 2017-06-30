package objects

import (
	"../../lib/es"
	"../locate"
	"io"
	"log"
	"net/http"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.Path, "/")[2]
	_, hash, e := es.SearchLatestVersion(name)
	if hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	s := locate.Locate(hash)
	if s == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	request, e := http.NewRequest("GET", "http://"+s+"/objects/"+hash, r.Body)
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
