package objects

import (
	"../../lib/es"
	"net/http"
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
	return digest[8:]
}

func getSizeFromHeader(r *http.Request) int64 {
	size, e := strconv.ParseInt(r.Header.Get("content-length"), 0, 64)
	if e != nil {
		panic(e)
	}
	return size
}

func addVersion(r *http.Request) error {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	hash := getHashFromHeader(r)
	size := getSizeFromHeader(r)

	version, _, e := es.SearchLatestVersion(name)
	if e != nil {
		return e
	}
	version += 1

	return es.PutVersion(name, version, size, hash)
}
