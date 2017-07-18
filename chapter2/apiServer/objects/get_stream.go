package objects

import (
	"../locate"
	"errors"
	"io"
	"lib/objectstream"
)

func getStream(object string) (io.Reader, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, errors.New("object locate fail")
	}
	return objectstream.NewGetStream(server, object)
}
