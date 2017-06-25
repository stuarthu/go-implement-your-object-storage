package cfg

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	RabbitMQServer string
	ListenAddress  string
	StorageRoot    string
}

var Config config

func init() {
	b, e := ioutil.ReadFile("config.json")
	if e != nil {
		panic(e)
	}
	if e = json.Unmarshal(b, &Config); e != nil {
		panic(e)
	}
}
