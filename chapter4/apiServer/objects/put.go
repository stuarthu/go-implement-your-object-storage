package objects

import (
	"../../lib/objectstream"
	"../heartbeat"
	"../locate"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func put(w http.ResponseWriter, r *http.Request) {
	hash := getHashFromHeader(r)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	info := locate.Locate(url.PathEscape(hash))
	if len(info) == 0 {
		c, e := storeObject(r)
		if e != nil {
			log.Println(e)
			w.WriteHeader(c)
			return
		}
		if c != http.StatusOK {
			w.WriteHeader(c)
			return
		}
	}

	e := addVersion(r)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func storeObject(r *http.Request) (int, error) {
	s := heartbeat.ChooseRandomDataServer()
	if s == "" {
		return http.StatusServiceUnavailable, fmt.Errorf("cannot find any dataServer")
	}
	hash := getHashFromHeader(r)
	size := getSizeFromHeader(r)
	h := sha256.New()
	reader := io.TeeReader(r.Body, h)
	stream, e := objectstream.NewTempStream(s, url.PathEscape(hash), size)
	if e != nil {
		return http.StatusInternalServerError, e
	}
	_, e = io.Copy(stream, reader)
	if e != nil {
		stream.Close(false)
		return http.StatusInternalServerError, e
	}
	digest := base64.StdEncoding.EncodeToString(h.Sum(nil))
	if digest != hash {
		stream.Close(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calculated=%s, requested=%s", digest, hash)
	}
	return stream.Close(true)
}
