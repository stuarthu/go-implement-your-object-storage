package versions

import (
	"../../lib/es"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	code, body, e := es.SearchVersions(strings.Split(r.URL.EscapedPath(), "/")[2])
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	dec := json.NewDecoder(body)
	enc := json.NewEncoder(w)
	for {
		t, e := dec.Token()
		if e != nil {
			if e != io.EOF {
				log.Println(e)
			}
			return
		}
		if t == "_source" {
			var m es.Metadata
			e = dec.Decode(&m)
			if e != nil {
				log.Println(e)
			}
			e = enc.Encode(&m)
			if e != nil {
				log.Println(e)
			}
		}
	}
}
