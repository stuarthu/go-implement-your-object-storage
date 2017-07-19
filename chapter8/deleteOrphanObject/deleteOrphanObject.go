package main

import (
	"../apiServer/heartbeat"
	"bufio"
	"lib/es"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	go heartbeat.ListenHeartbeat()
	time.Sleep(10 * time.Second)

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
			hash := strings.Split(object, ".")[0]
			objectInMetadata, e := es.HasObject(hash)
			if e != nil {
				log.Println(e)
				return
			}
			if !objectInMetadata {
				del(dataServer, object)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func del(server, object string) {
	log.Println("delete", object)
	url := "http://" + server + "/objects/" + object
	request, _ := http.NewRequest("DELETE", url, nil)
	client := http.Client{}
	client.Do(request)
}
