package temp

import (
	"io"
	"lib/objectstream"
	"log"
	"net/http"
	"strings"
)

func patch(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.EscapedPath(), "/")[2]
	t, e := tokenFromString(s)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	stream := &objectstream.TempPutStream{t.Server, t.Uuid}
	io.Copy(stream, r.Body)
}
