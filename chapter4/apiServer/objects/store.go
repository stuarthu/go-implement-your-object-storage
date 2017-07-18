package objects

import (
	"../locate"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func StoreObject(r io.Reader, hash string, size int64) (int, error) {
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}

	stream, e := putStream(url.PathEscape(hash), size)
	if e != nil {
		return http.StatusInternalServerError, e
	}

	h := sha256.New()
	reader := io.TeeReader(r, h)
	io.Copy(stream, reader)

	digest := base64.StdEncoding.EncodeToString(h.Sum(nil))
	if digest != hash {
		stream.Commit(false)
		return http.StatusForbidden, fmt.Errorf("object hash mismatch, calculated=%s, requested=%s", digest, hash)
	}
	stream.Commit(true)
	return http.StatusOK, nil
}
