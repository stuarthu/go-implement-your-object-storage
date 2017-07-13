package temp

import (
	"compress/gzip"
	"io"
	"os"
)

func copyTempToObjects(f *os.File, name string) {
	w, _ := os.Create(name)
	w2 := gzip.NewWriter(w)
	io.Copy(w2, f)
	w2.Close()
}
