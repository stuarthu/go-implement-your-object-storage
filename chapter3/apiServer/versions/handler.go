package versions

import (
	"../../lib/es"
	"io"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	code, body, e := es.SearchVersions(strings.Split(r.URL.Path, "/")[2])
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	io.Copy(w, body)
}
