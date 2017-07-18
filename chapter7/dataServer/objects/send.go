package objects

import (
	"compress/gzip"
	"io"
	"os"
)

func sendObject(w io.Writer, object string) {
	f, e := os.Open(object)
	if e != nil {
		panic(e)
	}
	defer f.Close()
	gzipStream, e := gzip.NewReader(f)
	defer gzipStream.Close()
	io.Copy(w, gzipStream)
}
