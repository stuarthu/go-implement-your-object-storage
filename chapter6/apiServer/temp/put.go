package temp

import (
	"../objects"
	"lib/objectstream"
	"lib/utils"
	"log"
	"net/http"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.EscapedPath(), "/")[2]
	t, e := tokenFromString(s)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	putStream := &objectstream.TempPutStream{t.Server, t.Uuid}
	defer putStream.Close(false)
	getStream, e := putStream.NewTempGetStream()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	c, e := objects.StoreObject(getStream, t.Hash, t.Size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(c)
		return
	}
	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}
	e = utils.AddVersion(t.Name, t.Hash, t.Size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
