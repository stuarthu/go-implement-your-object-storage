package objects

import (
	"../locate"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	filename := os.Getenv("STORAGE_ROOT") + "/" + name
	if locate.Locate(filename) {
		return
	}

	hash, e := url.PathUnescape(name)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	output, _ := exec.Command("uuidgen").Output()
	uuid := strings.TrimSuffix(string(output), "\n")
	tmpname := os.Getenv("STORAGE_ROOT") + "/tmp:" + uuid
	f, e := os.Create(tmpname)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reader := io.TeeReader(r.Body, f)
	h := sha256.New()
	io.Copy(h, reader)
	f.Close()
	digest := base64.StdEncoding.EncodeToString(h.Sum(nil))
	if digest != hash {
		os.Remove(tmpname)
		log.Println("calculated digest=" + digest + ",requested object=" + hash)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	os.Rename(tmpname, filename)
}
