package objectstream

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func TestGet(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(getHandler))
	defer s.Close()

	gs, _ := NewGetStream(s.URL[7:], "any")
	b, _ := ioutil.ReadAll(gs)
	if string(b) != "hello world" {
		t.Error(b)
	}
}
