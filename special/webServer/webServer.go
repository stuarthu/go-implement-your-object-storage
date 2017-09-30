package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	http.HandleFunc("/", webHandler)
	http.HandleFunc("/upload", uploadHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

type Metadata struct {
	Name    string
	Version int
	Size    int64
	Hash    string
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	req, e := http.Get("http://10.29.2.11:12345/versions/")
	if e != nil {
		log.Println(e)
		return
	}
	s := bufio.NewScanner(req.Body)

	w.Write([]byte("<html><table><tr><th>文件名</th><th>版本</th><th>大小</th></tr>"))
	for s.Scan() {
		var meta Metadata
		json.Unmarshal([]byte(s.Text()), &meta)
		if meta.Hash != "" {
			n, _ := url.PathUnescape(meta.Name)
			l := fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%d</td></tr>", n, meta.Version, meta.Size)
			w.Write([]byte(l))
		}
	}
	w.Write([]byte("</table>"))
	w.Write([]byte(`<form action=/upload method=post enctype=multipart/form-data><input type=file name=upload><input type=submit></form>`))
	w.Write([]byte("</html>"))
}

type Sizer interface {
	Size() int64
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	f, header, e := r.FormFile("upload")
	if e != nil {
		log.Println(e)
		return
	}
	defer f.Close()
	h := sha256.New()
	io.Copy(h, f)
	d := base64.StdEncoding.EncodeToString(h.Sum(nil))
	log.Println(d)
	f.Seek(0, 0)
	dat, _ := ioutil.ReadAll(f)
	req, e := http.NewRequest("PUT", "http://10.29.2.11:12345/objects/"+url.PathEscape(header.Filename), bytes.NewBuffer(dat))
	if e != nil {
		log.Println(e)
		return
	}
	req.Header.Set("digest", "SHA-256="+d)
	client := http.Client{}
	log.Println("uploading file", header.Filename, "hash", d, "size", f.(Sizer).Size())
	_, e = client.Do(req)
	if e != nil {
		log.Println(e)
		return
	}
	log.Println("uploaded")
	time.Sleep(time.Second)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
