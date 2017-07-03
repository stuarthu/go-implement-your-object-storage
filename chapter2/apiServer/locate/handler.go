package locate

import (
	"../../lib/rabbitmq"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Locate(name string) string {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	msg := <-c
	if len(msg.Body) == 0 {
		return ""
	}
	s, e := strconv.Unquote(string(msg.Body))
	if e != nil {
		panic(e)
	}
	return s
}

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	s := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
	if s != "" {
		w.Write([]byte(s))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
