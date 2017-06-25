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
	return !os.IsNotExist(err)
}

func StartLocate() {
	q := rabbitmq.New(cfg.Config.RabbitMQServer)
	defer q.Close()
	q.Bind("locate")
	c := q.Consume()
	var body interface{}
	for msg := range c {
		e := json.Unmarshal(msg.Body, &body)
		if e != nil {
			panic(e)
		}
		if locate(cfg.Config.StorageRoot + "/" + body.(string)) {
			fmt.Println(body, "exist, reply", cfg.Config.ListenAddress)
			q.Send(msg.ReplyTo, cfg.Config.ListenAddress)
		} else {
			fmt.Println(body, "not exist, reply to", msg.ReplyTo)
		}
	}
}
