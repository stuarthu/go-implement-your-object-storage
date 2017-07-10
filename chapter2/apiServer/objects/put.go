package objects

import (
	"../../lib/httpstream"
	"../heartbeat"
	"io"
	"log"
	"net/http"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	s := heartbeat.ChooseRandomDataServer()
	if s == "" {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream := httpstream.NewPutStream("http://" + s + "/objects/" + object)
	io.Copy(stream, r.Body)
	e := stream.Close()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
