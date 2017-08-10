package rs

import (
	"io"
	"lib/objectstream"
)

type RSResumableGetStream struct {
	*decoder
}

func NewRSResumableGetStream(dataServers []string, uuids []string, size int64) (*RSResumableGetStream, error) {
	readers := make([]io.Reader, ALL_SHARDS)
	var e error
	for i := 0; i < ALL_SHARDS; i++ {
		readers[i], e = objectstream.NewTempGetStream(dataServers[i], uuids[i])
		if e != nil {
			return nil, e
		}
	}
	writers := make([]io.Writer, ALL_SHARDS)
	dec := NewDecoder(readers, writers, size)
	return &RSResumableGetStream{dec}, nil
}
