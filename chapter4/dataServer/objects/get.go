package objects

import (
	"net/http"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	object := getObject(strings.Split(r.URL.EscapedPath(), "/")[2])
	if object == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	sendObject(w, object)
}
