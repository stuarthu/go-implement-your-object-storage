package temp

import (
	"fmt"
	"lib/rs"
	"log"
	"net/http"
	"strings"
)

func head(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, e := rs.NewRSResumablePutStreamFromToken(s)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("content-length", fmt.Sprintf("%d", current))
}
