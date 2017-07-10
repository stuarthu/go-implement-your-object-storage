package objects

import (
	"../../lib/es"
	"../../lib/httpstream"
	"../heartbeat"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func getHashFromHeader(r *http.Request) string {
	digest := r.Header.Get("digest")
	if len(digest) < 9 {
		return ""
	}
	if digest[:8] != "SHA-256=" {
		return ""
	}
	return url.PathEscape(digest[8:])
}

func addVersion(r *http.Request, hash string) error {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]

	size, e := strconv.ParseInt(r.Header.Get("content-length"), 0, 64)
	if e != nil {
		return e
	}

	version, _, e := es.SearchLatestVersion(name)
	if e != nil {
		return e
	}
	version += 1

	return es.PutVersion(name, version, size, hash)
}

func storeObject(r *http.Request, hash string) error {
	s := heartbeat.ChooseRandomDataServer()
	if s == "" {
		return fmt.Errorf("cannot find any dataServer")
	}
	stream := httpstream.NewPutStream("http://" + s + "/objects/" + hash)
	io.Copy(stream, r.Body)
	return stream.Close()
}
