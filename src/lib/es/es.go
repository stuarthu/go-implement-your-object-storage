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

type Metadata struct {
	Name    string
	Version int
	Size    int64
	Hash    string
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total int
		Hits  []hit
	}
}

func GetVersion(name string, versionId int) (version Metadata, e error) {
	client := http.Client{}
	request, _ := http.NewRequest("GET",
		fmt.Sprintf("http://%s/metadata/objects/%s_%d/_source?",
			os.Getenv("ES_SERVER"), name, version), nil)
	r, e := client.Do(request)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to get %s_%d: %d", name, version, r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	e = json.Unmarshal(result, &version)
	return
}

func SearchLatestVersion(name string) (version Metadata, e error) {
	client := http.Client{}
	request, _ := http.NewRequest("GET", "http://"+os.Getenv("ES_SERVER")+
		"/metadata/objects/_search?q=name:"+name+
		"&size=1&sort=version:desc", nil)
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
	if sr.Hits.Total != 0 {
		version = sr.Hits.Hits[0].Source
	}
	return
}

func PutVersion(name string, version int, size int64, hash string) error {
	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s"}`,
		name, version, size, hash)
	client := http.Client{}
	request, _ := http.NewRequest("PUT",
		fmt.Sprintf("http://%s/metadata/objects/%s_%d?op_type=create",
			os.Getenv("ES_SERVER"), name, version),
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

func SearchVersions(name string) (int, io.Reader, error) {
	client := http.Client{}
	url := "http://" + os.Getenv("ES_SERVER") + "/metadata/objects/_search"
	if name != "" {
		url += "?q=name:" + name
	}
	request, _ := http.NewRequest("GET", url, nil)
	r, e := client.Do(request)
	if e != nil {
		return 0, nil, e
	}
	return r.StatusCode, r.Body, nil
}
