package objects

import (
	"io"
	"os"
)

func sendObject(w io.Writer, object string) {
	f, _ := os.Open(object)
	defer f.Close()
	io.Copy(w, f)
}
