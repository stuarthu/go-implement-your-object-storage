package objects

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" +
		strings.Split(r.URL.EscapedPath(), "/")[2] + ".*")
	if len(files) != 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	object := files[0]
	hash := strings.Split(object, ".")[2]
	f, _ := os.Open(object)
	defer f.Close()
	h := sha256.New()
	io.Copy(h, f)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	if d != hash {
		log.Println("object hash mismatch, remove", object)
		os.Remove(object)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	f.Seek(0, 0)
	io.Copy(w, f)
}
