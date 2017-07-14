package objects

import (
	"../heartbeat"
	"../locate"
	"errors"
	"io"
	"lib/es"
	"lib/rs"
	"net/url"
)

func getData(meta es.Metadata) (io.Reader, error) {
	object := url.PathEscape(meta.Hash)
	locateInfo := locate.Locate(object)
	if len(locateInfo) < rs.DATA_SHARDS {
		return nil, errors.New("object locate fail")
	}
	ds := make([]string, 0)
	if len(locateInfo) != rs.ALL_SHARDS {
		ds = heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	return rs.NewRSGetStream(locateInfo, ds, object, meta.Size)
}
