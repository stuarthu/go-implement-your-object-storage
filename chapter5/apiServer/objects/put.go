package objects

import (
	"../../lib/es"
	"../../lib/rs"
	"../heartbeat"
	"../locate"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func storeData(w http.ResponseWriter, r *http.Request, object string, size int64) {
	s := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS)
	if len(s) != rs.ALL_SHARDS {
		log.Println("cannot find enough dataServer")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	log.Println("rs.Put", object, size)
	e := rs.Put(r.Body, s, object, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
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

	size, e := strconv.ParseInt(r.Header.Get("content-length"), 0, 64)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	object := url.PathEscape(hash)
	s := locate.Locate(object)
	if len(s) == 0 {
		storeData(w, r, object, size)
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	version, _, e := es.SearchLatestVersion(name)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	version += 1
	e = es.PutVersion(name, version, size, hash)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
