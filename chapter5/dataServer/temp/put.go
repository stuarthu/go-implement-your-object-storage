package temp

import (
	"../locate"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, e := readFromFile(uuid)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	datFile := infoFile + ".dat"
	f, e := os.OpenFile(datFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	os.Remove(infoFile)
	if actual != tempinfo.Size {
		os.Remove(datFile)
		log.Println("actual size mismatch")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h := sha256.New()
	f2, _ := os.Open(datFile)
	defer f2.Close()
	io.Copy(h, f2)
	d := base64.StdEncoding.EncodeToString(h.Sum(nil))
	locate.Add(tempinfo.object(), tempinfo.id())
	os.Rename(datFile, os.Getenv("STORAGE_ROOT")+"/objects/"+
		tempinfo.Name+"."+url.PathEscape(d))
}
