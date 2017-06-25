package locate

import (
	"../../lib/cfg"
	"../../lib/rabbitmq"
	"encoding/json"
	"fmt"
	"os"
)

func locate(name string) bool {
	_, err := os.Stat(name)
	return os.IsExist(err)
}

func ListenLocate() {
	q := rabbitmq.New(cfg.C.RabbitMQServer, "locate", "")
	defer q.Close()
	c := q.Consume()
	var body interface{}
	for msg := range c {
		e := json.Unmarshal(msg.Body, &body)
		if e != nil {
			panic(e)
		}
		if locate(cfg.C.StorageRoot + body.(string)) {
			fmt.Println(body, " exist")
		} else {
			fmt.Println(body, " not exist")
		}
	}
}
