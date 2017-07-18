package objects

import (
	"../locate"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func getObject(name string) string {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + name + ".*")
	if len(files) != 1 {
		return ""
	}
	object := files[0]
	h := sha256.New()
	sendObject(h, object)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	hash := strings.Split(object, ".")[2]
	if d != hash {
		log.Println("object hash mismatch, remove", object)
		locate.Del(name)
		os.Remove(object)
		return ""
	}
	return object
}
