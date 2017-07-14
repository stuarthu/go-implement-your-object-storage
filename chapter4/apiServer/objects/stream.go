package objects

import (
	"lib/objectstream"
	"../heartbeat"
	"net/url"
    "fmt"
)

func createStream(hash string, size int64) (*objectstream.TempPutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}

	return objectstream.NewTempPutStream(server, url.PathEscape(hash), size)
}
