package heartbeat

import (
	"../../lib/rabbitmq"
	"log"
	"math/rand"
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

func ChooseRandomDataServers(s int) (ds []string) {
	good := make([]string, 0)
	bad := make([]string, 0)
	for s, t := range DataServers {
		if t.Add(30 * time.Second).After(time.Now()) {
			good = append(good, s)
		} else {
			bad = append(bad, s)
		}
	}
	log.Println(good, bad)
	if len(good) < s {
		return
	}
	p := rand.Perm(s)
	for i := range p {
		ds = append(ds, good[i])
	}
	mutex.Lock()
	for i := range bad {
		delete(DataServers, bad[i])
	}
	mutex.Unlock()
	return
}
