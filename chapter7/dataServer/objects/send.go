package objects

import (
	"compress/gzip"
	"io"
	"os"
)

func sendObject(w io.Writer, object string) {
	f, _ := os.Open(object)
	defer f.Close()
	gzipStream, _ := gzip.NewReader(f)
	defer gzipStream.Close()
	io.Copy(w, gzipStream)
}
