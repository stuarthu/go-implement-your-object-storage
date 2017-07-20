package heartbeat

import (
	"lib/rabbitmq"
	"os"
	"strconv"
	"sync"
	"time"
)

var DataServers = make(map[string]time.Time)
var mutex sync.Mutex

func ListenHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("apiServers")
	c := q.Consume()
	go removeExpiredDataServer()
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		DataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		for s, t := range DataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				mutex.Lock()
				delete(DataServers, s)
				mutex.Unlock()
			}
		}
	}
}
