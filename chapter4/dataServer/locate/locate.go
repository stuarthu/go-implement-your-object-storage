package locate

import (
	"../../lib/rabbitmq"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type locateMessage struct {
	Addr string
	Id   int
}

var objects = make(map[string]int)
var mutex sync.Mutex

func Locate(object string) int {
	mutex.Lock()
	id, _ := objects[object]
	mutex.Unlock()
	return id
}

func Add(object string, id int) {
	mutex.Lock()
	objects[object] = id
	mutex.Unlock()
}

func StartLocate() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		object, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		id := Locate(object)
		if id != 0 {
			q.Send(msg.ReplyTo, locateMessage{Addr: os.Getenv("LISTEN_ADDRESS"), Id: id})
		}
	}
}

func CollectObjects() {
	files, e := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	if e != nil {
		panic(e)
	}
	for i := range files {
		object := filepath.Base(files[i])
		objects[object] = 1
	}
}
