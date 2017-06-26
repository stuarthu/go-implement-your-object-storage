package objects

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func Put(w http.ResponseWriter, r *http.Request) {
	f, e := os.Create("/tmp/" + strings.Split(r.URL.Path, "/")[2])
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, r.Body)
}
