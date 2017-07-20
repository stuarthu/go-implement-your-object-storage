package objects

import (
	"compress/gzip"
	"io"
	"log"
	"os"
)

func sendObject(w io.Writer, object string) {
	f, e := os.Open(object)
	if e != nil {
		log.Println(e)
		return
	}
	defer f.Close()
	gzipStream, e := gzip.NewReader(f)
	if e != nil {
		log.Println(e)
		return
	}
	io.Copy(w, gzipStream)
	gzipStream.Close()
}
