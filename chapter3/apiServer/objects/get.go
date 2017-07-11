package objects

import (
	"../../lib/es"
	"../../lib/objectstream"
	"../locate"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	version := r.URL.Query()["version"]
	var hash string
	if len(version) == 0 {
		_, hash, _ = es.SearchLatestVersion(name)
	} else {
		v, e := strconv.Atoi(version[0])
		if e != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		hash, _ = es.GetHash(name, v)
	}
	if hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	object := url.PathEscape(hash)
	info := locate.Locate(object)
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	stream, e := objectstream.NewGetStream(info[0].Addr, object)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(w, stream)
}
