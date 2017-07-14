package objects

import (
	"lib/objectstream"
    "io"
    "net/url"
	"../locate"
	"errors"
	"lib/es"
)

func getData(meta es.Metadata) (io.Reader, error) {
	object := url.PathEscape(meta.Hash)
	server := locate.Locate(object)
	if server == "" {
		return nil, errors.New("object locate fail")
	}
	return objectstream.NewGetStream(server, object)
}
