package es

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type hit struct {
	Source struct {
		Version int
		Hash    string
	} `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total int
		Hits  []hit
	}
}

func GetHash(name string, version int) (string, error) {
	client := http.Client{}
	request, _ := http.NewRequest("GET", "http://"+os.Getenv("ES_SERVER")+
		fmt.Sprintf("http://%s/metadata/objects/%s_%d/_source?_source_include=hash", os.Getenv("ES_SERVER"), name, version), nil)
	r, e := client.Do(request)
	if e != nil {
		return "", e
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to get %s_%d: %d", name, version, r.StatusCode)
		return "", e
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	fmt.Println(result)
	e = json.Unmarshal(result, &sr)
	if e != nil {
		return "", e
	}
	return "", nil
}

func SearchLatestVersion(name string) (version int, hash string, e error) {
	client := http.Client{}
	request, _ := http.NewRequest("GET", "http://"+os.Getenv("ES_SERVER")+
		"/metadata/objects/_search?q=name:"+name+
		"&size=1&sort=version:desc&_source_include=version,hash", nil)
	r, e := client.Do(request)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search latest version: %d", r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	e = json.Unmarshal(result, &sr)
	if e != nil {
		return
	}
	if sr.Hits.Total == 0 {
		version = 0
		return
	}
	version = sr.Hits.Hits[0].Source.Version
	hash = sr.Hits.Hits[0].Source.Hash
	return
}

func PutVersion(name string, version, size int, hash string) error {
	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s"}`, name, version, size, hash)
	client := http.Client{}
	request, _ := http.NewRequest("PUT",
		fmt.Sprintf("http://%s/metadata/objects/%s_%d?op_type=create", os.Getenv("ES_SERVER"), name, version),
		strings.NewReader(doc))
	r, e := client.Do(request)
	if e != nil {
		return e
	}
	if r.StatusCode == http.StatusConflict {
		return PutVersion(name, version+1, size, hash)
	}
	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("fail to put version: %d %s", r.StatusCode, string(result))
	}
	return nil
}

func SearchVersions(object string) (int, io.Reader, error) {
	client := http.Client{}
	url := "http://" + os.Getenv("ES_SERVER") + "/metadata/objects/_search"
	if object != "" {
		url += "?q=name:" + object
	}
	request, _ := http.NewRequest("GET", url, nil)
	r, e := client.Do(request)
	if e != nil {
		return 0, nil, e
	}
	return r.StatusCode, r.Body, nil
}
