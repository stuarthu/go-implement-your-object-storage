package objects

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		Put(w, r)
		return
	}
	if m == http.MethodGet {
		Get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
