package objects

import (
	"io"
	"net/http"
	"net/url"
)

func StoreObject(r io.Reader, hash string) (int, error) {
	stream, e := createStream(url.PathEscape(hash))
	if e != nil {
		return http.StatusInternalServerError, e
	}

	io.Copy(stream, r)
	e = stream.Close()
	if e != nil {
		return http.StatusInternalServerError, e
	}
	return http.StatusOK, nil
}
