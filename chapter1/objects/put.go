package objects

import (
	"io"
	"net/http"
	"os"
)

func Put(w http.ResponseWriter, r *http.Request) {
	f, e := os.Create("/tmp" + r.URL.Path)
	defer f.Close()
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(f, r.Body)
}
