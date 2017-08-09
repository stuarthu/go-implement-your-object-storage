package temp

import (
	"../objects"
    "lib/es"
	"lib/objectstream"
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
	tempPutStream := &objectstream.TempPutStream{t.Server, t.Uuid}
	defer tempPutStream.Commit(false)
	tempGetStream, e := tempPutStream.NewTempGetStream()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	c, e := objects.StoreObject(tempGetStream, t.Hash, t.Size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(c)
		return
	}
	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}
	e = es.AddVersion(t.Name, t.Hash, t.Size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
