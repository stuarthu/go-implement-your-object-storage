package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type GetStream struct {
	reader io.Reader
}

func NewGetStream(server, object string) (*GetStream, error) {
	r, e := http.Get("http://" + server + "/objects/" + object)
	if e != nil {
		return nil, e
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code %d", r.StatusCode)
	}
	return &GetStream{r.Body}, nil
}

func (r *GetStream) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}
