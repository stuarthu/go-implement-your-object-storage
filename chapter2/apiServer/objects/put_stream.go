package objects

import (
	"../heartbeat"
	"fmt"
	"lib/objectstream"
	"net/url"
)

func putStream(object string) (*objectstream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}

	return objectstream.NewPutStream(server, url.PathEscape(object)), nil
}
