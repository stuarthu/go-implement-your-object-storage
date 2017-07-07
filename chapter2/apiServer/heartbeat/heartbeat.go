package heartbeat

import (
	"../../lib/rabbitmq"
	"math/rand"
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
		mutex.Unlock()
	}
}

func ChooseRandomDataServer() string {
	n := len(dataServers)
	if n == 0 {
		return ""
	}
	i := rand.Intn(n)
	for s, t := range dataServers {
		if i == 0 {
			if t.Add(30 * time.Second).After(time.Now()) {
				return s
			}
			mutex.Lock()
			delete(dataServers, s)
			mutex.Unlock()
			return ChooseRandomDataServer()
		}
		i--
	}
	panic("should not pass here")
}
