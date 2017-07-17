package objects

import (
	"../heartbeat"
	"fmt"
	"lib/objectstream"
	"net/url"
)

func createStream(hash string) (*objectstream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}

	return objectstream.NewPutStream(server, url.PathEscape(hash)), nil
}
