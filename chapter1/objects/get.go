package objects

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func Get(w http.ResponseWriter, r *http.Request) {
	f, e := os.Open("/tmp/" + strings.Split(r.URL.Path, "/")[2])
	defer f.Close()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
	}
	io.Copy(w, f)
}
