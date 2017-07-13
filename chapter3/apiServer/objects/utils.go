package objects

import (
	"lib/es"
	"net/http"
	"strconv"
	"strings"
)

func getSizeFromHeader(r *http.Request) int64 {
	size, e := strconv.ParseInt(r.Header.Get("content-length"), 0, 64)
	if e != nil {
		panic(e)
	}
	return size
}
