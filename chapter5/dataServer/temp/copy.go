package temp

import (
    "../locate"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/url"
    "os"
)

func copyTempToObjects(f *os.File, tempinfo *tempInfo) {
	h := sha256.New()
	io.Copy(h, f)
	d := base64.StdEncoding.EncodeToString(h.Sum(nil))
	locate.Add(tempinfo.object(), tempinfo.id())
    w, _ := os.Create(os.Getenv("STORAGE_ROOT")+"/objects/"+tempinfo.Name+"."+url.PathEscape(d))
    f.Seek(0, 0)
    io.Copy(w, f)
}
