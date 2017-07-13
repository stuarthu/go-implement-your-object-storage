package objects

import (
	"../locate"
	"crypto/sha256"
	"encoding/base64"
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
	h := sha256.New()
	sendObject(h, object)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	hash := strings.Split(object, ".")[2]
	if d != hash {
		log.Println("object hash mismatch, remove", object)
		locate.Del(object)
		os.Remove(object)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sendObject(w, object)
}
