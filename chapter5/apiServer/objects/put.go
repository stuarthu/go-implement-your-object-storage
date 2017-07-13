package objects

import (
	"../heartbeat"
	"../locate"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"lib/rs"
	"lib/utils"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	hash := utils.GetHashFromHeader(r)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	info := locate.Locate(url.PathEscape(hash))
	size := utils.GetSizeFromHeader(r)
	if len(info) == 0 {
		c, e := StoreObject(r.Body, hash, size)
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

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	e := utils.AddVersion(name, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func StoreObject(r io.Reader, hash string, size int64) (int, error) {
	ds := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(ds) < rs.ALL_SHARDS {
		return http.StatusServiceUnavailable, fmt.Errorf("cannot find enough dataServer")
	}
	h := sha256.New()
	reader := io.TeeReader(r, h)
	stream, e := rs.NewRSPutStream(ds, url.PathEscape(hash), size)
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
