package objects

import (
	"../../lib/es"
	"../heartbeat"
	"../locate"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type searchResult struct {
	Hits struct {
		Total int
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	digest := r.Header.Get("digest")
	if len(digest) < 9 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if digest[:8] != "SHA-256=" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hash := digest[8:]

	object := url.PathEscape(hash)
	s := locate.Locate(object)
	if s != "" {
		return
	}

	s = heartbeat.ChooseRandomDataServer()
	if s == "" {
		log.Println("cannot find any dataServer")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

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

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	version, _, e := es.SearchLatestVersion(name)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	version += 1
	size, e := strconv.Atoi(r.Header.Get("content-length"))
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	e = es.PutVersion(name, version, size, hash)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
