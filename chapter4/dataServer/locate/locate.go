package locate

import (
	"lib/rabbitmq"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var objects = make(map[string]int)
var mutex sync.Mutex

func Locate(object string) bool {
	mutex.Lock()
	_, ok := objects[object]
	mutex.Unlock()
	return ok
}

func Add(object string) {
	mutex.Lock()
	objects[object] = 1
	mutex.Unlock()
}

func Del(object string) {
	mutex.Lock()
	delete(objects, object)
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
		exist := Locate(object)
		if exist {
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}

func CollectObjects() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		object := filepath.Base(files[i])
		objects[object] = 1
	}
}
