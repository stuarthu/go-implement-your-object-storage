package locate

import (
	"../../lib/rabbitmq"
	"os"
	"strconv"
)

func locate(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func StartLocate() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		n, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		if locate(os.Getenv("STORAGE_ROOT") + "/" + n) {
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}
