package heartbeat

import (
	"lib/rabbitmq"
	"os"
	"strconv"
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

func ListenHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("apiServers")
	c := q.Consume()
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		removeExpiredDataServer()
		mutex.Unlock()
	}
}

func removeExpiredDataServer() {
	for s, t := range dataServers {
		if t.Add(30 * time.Second).Before(time.Now()) {
			delete(dataServers, s)
		}
	}
}
