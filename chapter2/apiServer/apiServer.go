package main

import (
	"../lib/cfg"
	"../lib/rabbitmq"
	"./objects"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func objectsHandler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		objects.Put(w, r)
		return
	}
	if m == http.MethodGet {
		objects.Get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func locateHandler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	url := strings.Split(r.URL.Path, "/")
	q := rabbitmq.New(cfg.Config.RabbitMQServer)
	defer q.Close()
	q.Publish("locate", url[2])
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Send(q.Name, "not found")
		fmt.Println("timeout sent to", q.Name)
	}()
	msg := <-c
	w.Write(msg.Body)
}

func main() {
	http.HandleFunc("/objects/", objectsHandler)
	http.HandleFunc("/locate/", locateHandler)
	log.Fatal(http.ListenAndServe(cfg.Config.ListenAddress, nil))
}
