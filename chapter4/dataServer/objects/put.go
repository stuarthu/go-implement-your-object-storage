package objects

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]

	files, _ := filepath.Glob(os.Getenv("TMP_ROOT") + "/" + object + ":*")
	if len(files) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(files) > 1 {
		w.WriteHeader(http.StatusConflict)
		return
	}
	os.Rename(files[0], os.Getenv("STORAGE_ROOT")+"/"+object)
}
