package versions

import (
	"../../lib/es"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	es.SearchVersions(w, strings.Split(r.URL.Path, "/")[2])
	return
}
