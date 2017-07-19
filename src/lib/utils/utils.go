package utils

import (
	"lib/es"
	"net/http"
	"strconv"
)

func GetHashFromHeader(r *http.Request) string {
	digest := r.Header.Get("digest")
	if len(digest) < 9 {
		return ""
	}
	if digest[:8] != "SHA-256=" {
		return ""
	}
	return digest[8:]
}

func AddVersion(name, hash string, size int64) error {
	version, e := es.SearchLatestVersion(name)
	if e != nil {
		return e
	}
	return es.PutMetadata(name, version.Version+1, size, hash)
}

func GetSizeFromHeader(r *http.Request) int64 {
	size, _ := strconv.ParseInt(r.Header.Get("content-length"), 0, 64)
	return size
}
