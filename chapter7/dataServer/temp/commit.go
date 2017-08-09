package temp

import (
	"../locate"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func (t *tempInfo) hash() string {
	s := strings.Split(t.Hash, ".")
	return s[0]
}

func (t *tempInfo) id() int {
	s := strings.Split(t.Hash, ".")
	id, _ := strconv.Atoi(s[1])
	return id
}

func commitTempObject(datFile string, tempinfo *tempInfo) {
	f, _ := os.Open(datFile)
	defer f.Close()
	h := sha256.New()
	io.Copy(h, f)
	d := base64.StdEncoding.EncodeToString(h.Sum(nil))
	f.Seek(0, 0)
	w, _ := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + tempinfo.Hash + "." + url.PathEscape(d))
	w2 := gzip.NewWriter(w)
	io.Copy(w2, f)
	w2.Close()
	os.Remove(datFile)
	locate.Add(tempinfo.hash(), tempinfo.id())
}
