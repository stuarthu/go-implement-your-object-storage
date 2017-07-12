package heartbeat

import (
	"lib/rabbitmq"
	"os"
	"time"
)

func StartHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	for {
		time.Sleep(5 * time.Second)
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
	}
}
