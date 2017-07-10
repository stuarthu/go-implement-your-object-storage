package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type PutStream struct {
	writer *io.PipeWriter
	c      chan error
}

func NewPutStream(server, object string) *PutStream {
	reader, writer := io.Pipe()
	request, _ := http.NewRequest("PUT", "http://"+server+"/objects/"+object, reader)
	client := http.Client{}
	c := make(chan error)
	go func() {
		r, e := client.Do(request)
		if e == nil && r.StatusCode != http.StatusOK {
			e = fmt.Errorf("http code %d", r.StatusCode)
		}
		c <- e
	}()
	return &PutStream{writer, c}
}

func (w *PutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *PutStream) Close() error {
	w.writer.Close()
	return <-w.c
}
