package versions

import (
	"encoding/json"
	"io"
	"lib/es"
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
	more := true
	from := 0
	size := 1000
	for more {
		body, e := es.SearchVersions(strings.Split(r.URL.EscapedPath(), "/")[2], from, size)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		from += size
		dec := json.NewDecoder(body)
		enc := json.NewEncoder(w)
		count := 0
		for {
			t, e := dec.Token()
			if e != nil {
				if e != io.EOF {
					log.Println(e)
				}
                break
			}
			if t == "_source" {
				count++
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
		if count != size {
			more = false
		}
	}
}
