package objects

import (
	"io"
	"net/http"
	"net/url"
)

func StoreObject(r io.Reader, object string) (int, error) {
	stream, e := putStream(url.PathEscape(object))
	if e != nil {
		return http.StatusServiceUnavailable, e
	}

	io.Copy(stream, r)
	e = stream.Close()
	if e != nil {
		return http.StatusInternalServerError, e
	}
	return http.StatusOK, nil
}
