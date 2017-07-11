package objectstream

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type TempPutStream struct {
	server string
	uuid   string
}

func NewTempPutStream(server, object string, size int64) (*TempPutStream, error) {
	request, e := http.NewRequest("POST", "http://"+server+"/temp/"+object, nil)
	if e != nil {
		panic(e)
	}
	request.Header.Set("size", fmt.Sprintf("%d", size))
	client := http.Client{}
	response, e := client.Do(request)
	if e != nil {
		return nil, e
	}
	uuid, e := ioutil.ReadAll(response.Body)
	if e != nil {
		panic(e)
	}
	return &TempPutStream{server, string(uuid)}, nil
}

func (w *TempPutStream) Write(p []byte) (n int, err error) {
	request, e := http.NewRequest("PATCH", "http://"+w.server+"/temp/"+w.uuid, strings.NewReader(string(p)))
	if e != nil {
		return 0, e
	}
	client := http.Client{}
	r, e := client.Do(request)
	if e != nil {
		return 0, e
	}
	if r.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dataServer return http code %d", r.StatusCode)
	}
	return len(p), nil
}

func (w *TempPutStream) Close(good bool) (int, error) {
	method := "DELETE"
	if good {
		method = "PUT"
	}
	request, _ := http.NewRequest(method, "http://"+w.server+"/temp/"+w.uuid, nil)
	client := http.Client{}
	r, e := client.Do(request)
	return r.StatusCode, e
}

func (w *TempPutStream) NewTempGetStream() (*GetStream, error) {
	return newGetStream("http://" + w.server + "/temp/" + w.uuid)
}
