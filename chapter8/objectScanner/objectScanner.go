package main

import (
	"../apiServer/heartbeat"
	"../apiServer/objects"
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"lib/es"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	go heartbeat.ListenHeartbeat()
	time.Sleep(6 * time.Second)

	for dataServer, _ := range heartbeat.DataServers {
		url := "http://" + dataServer + "/objects/"
		r, e := http.Get(url)
		if e != nil {
			log.Println(e)
			return
		}
		scanner := bufio.NewScanner(r.Body)
		for scanner.Scan() {
			object := scanner.Text()
			tmp := strings.Split(object, ".")
			hash := tmp[0]
			id := tmp[1]
			if id == "0" {
				verify(hash)
			}
		}
	}
}

func verify(hash string) {
	size, e := es.SearchHashSize(hash)
	if e != nil {
		log.Println(e)
		return
	}
	stream, e := objects.GetStream(hash, size)
	if e != nil {
		log.Println(e)
		return
	}
	h := sha256.New()
	io.Copy(h, stream)
	d := base64.StdEncoding.EncodeToString(h.Sum(nil))
	if d != hash {
		log.Printf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}
}
