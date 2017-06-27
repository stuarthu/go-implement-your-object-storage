package objects

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func Get(w http.ResponseWriter, r *http.Request) {
	f, e := os.Open(os.Getenv("STORAGE_ROOT") + "/" + strings.Split(r.URL.Path, "/")[2])
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}
