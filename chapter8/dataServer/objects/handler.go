package objects

import (
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	if m == http.MethodGet {
		if name == "" {
			listObjects(w)
			return
		}
		get(w, r)
		return
	}
	if m == http.MethodDelete {
		del(name)
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
