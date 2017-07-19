package locate

import (
	"lib/rabbitmq"
	"lib/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var Objects = make(map[string]int)
var mutex sync.Mutex

func Locate(object string) int {
	mutex.Lock()
	id, ok := Objects[object]
	mutex.Unlock()
	if !ok {
		return -1
	}
	return id
}

func Add(object string, id int) {
	mutex.Lock()
	Objects[object] = id
	mutex.Unlock()
}

func Del(object string) {
	mutex.Lock()
	delete(Objects, object)
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
		if id != -1 {
			q.Send(msg.ReplyTo, types.LocateMessage{Addr: os.Getenv("LISTEN_ADDRESS"), Id: id})
		}
	}
}

func CollectObjects() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			panic(files[i])
		}
		object := file[0]
		id, e := strconv.Atoi(file[1])
		if e != nil {
			panic(e)
		}
		Objects[object] = id
	}
}
