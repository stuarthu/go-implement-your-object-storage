package objects

import (
	"../heartbeat"
	"../locate"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"lib/rs"
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
	s := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(s) < rs.ALL_SHARDS {
		return http.StatusServiceUnavailable, fmt.Errorf("cannot find enough dataServer")
	}
	hash := getHashFromHeader(r)
	size := getSizeFromHeader(r)
	h := sha256.New()
	reader := io.TeeReader(r.Body, h)
	stream, e := rs.NewRSPutStream(s, url.PathEscape(hash), size)
	if e != nil {
		return http.StatusInternalServerError, e
	}
	io.Copy(stream, reader)
	stream.Close()
	digest := base64.StdEncoding.EncodeToString(h.Sum(nil))
	if digest != hash {
		stream.Submit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calculated=%s, requested=%s", digest, hash)
	}
	stream.Submit(true)
	return http.StatusOK, nil
}
