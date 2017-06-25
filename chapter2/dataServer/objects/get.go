package objects

import (
	"io"
	"net/http"
	"os"
)

func Get(w http.ResponseWriter, r *http.Request) {
	f, e := os.Open("/tmp" + r.URL.Path)
	defer f.Close()
	if e != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	io.Copy(w, f)
}
