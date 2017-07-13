package objects

import (
	"../heartbeat"
	"../locate"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"lib/rs"
	"net/http"
	"net/url"
)

func StoreObject(r io.Reader, hash string, size int64) (int, error) {
	info := locate.Locate(url.PathEscape(hash))
	if len(info) != 0 {
		return http.StatusOK, nil
	}

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
