package temp

import (
	"../locate"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func (t *tempInfo) hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}

func (t *tempInfo) id() int {
	s := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(s[1])
	return id
}

func commitTempObject(datFile string, tempinfo *tempInfo) {
	f, _ := os.Open(datFile)
	defer f.Close()
	h := sha256.New()
	io.Copy(h, f)
	d := base64.StdEncoding.EncodeToString(h.Sum(nil))
	os.Rename(datFile, os.Getenv("STORAGE_ROOT")+"/objects/"+tempinfo.Name+"."+url.PathEscape(d))
	locate.Add(tempinfo.hash(), tempinfo.id())
}
