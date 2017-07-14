package objects

import (
	"../locate"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/url"
	"os"
)

func getObject(hash string) string {
	object := os.Getenv("STORAGE_ROOT") + "/objects/" + hash
	h := sha256.New()
	sendObject(h, object)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	if d != hash {
		log.Println("object hash mismatch, remove", object)
		locate.Del(hash)
		os.Remove(object)
		return ""
	}
    return object
}
