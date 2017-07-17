package objects

import (
	"../locate"
	"errors"
	"io"
	"lib/es"
	"lib/objectstream"
	"net/url"
)

func getData(meta es.Metadata) (io.Reader, error) {
	object := url.PathEscape(meta.Hash)
	server := locate.Locate(object)
	if server == "" {
		return nil, errors.New("object locate fail")
	}
	return objectstream.NewGetStream(server, object)
}
