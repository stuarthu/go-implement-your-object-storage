package objects

import (
	"../heartbeat"
	"../locate"
	"fmt"
	"lib/rs"
)

func GetStream(object string, size int64) (*rs.RSGetStream, error) {
	locateInfo := locate.Locate(object)
	if len(locateInfo) < rs.DATA_SHARDS {
		return nil, fmt.Errorf("object %s locate fail, result %v", object, locateInfo)
	}
	ds := make([]string, 0)
	if len(locateInfo) != rs.ALL_SHARDS {
		ds = heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	return rs.NewRSGetStream(locateInfo, ds, object, size)
}
