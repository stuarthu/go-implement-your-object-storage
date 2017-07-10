package locate

import (
	"../../lib/rabbitmq"
	"../../lib/rs"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
)

type locateMessage struct {
	Addr string
	Id   int
}

func Locate(name string) (locateInfo []locateMessage) {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	for i := 0; i < rs.ALL_SHARDS; i++ {
		msg := <-c
		if len(msg.Body) == 0 {
			return
		}
		var info locateMessage
		json.Unmarshal(msg.Body, &info)
		locateInfo = append(locateInfo, info)
	}
	return
}

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	for i := range info {
		b, _ := json.Marshal(info[i])
		w.Write(b)
		w.Write([]byte("\n"))
	}
}
