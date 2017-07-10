package objects

import (
	"../../lib/httpstream"
	"../locate"
	"io"
	"log"
	"net/http"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	s := locate.Locate(object)
	if s == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	stream, e := httpstream.NewGetStream("http://" + s + "/objects/" + object)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(w, stream)
}
