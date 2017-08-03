package rs

import (
	"fmt"
	"io"
	"lib/objectstream"
)

type RSPutStream struct {
	dataServers []string
	object      string
	size        int64
	writers     []io.Writer
	enc         *encoder
}

func NewRSPutStream(dataServers []string, object string, size int64) (*RSPutStream, error) {
	if len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	writers := make([]io.Writer, ALL_SHARDS)
	var e error
	for i := range writers {
		writers[i], e = objectstream.NewTempPutStream(dataServers[i],
			fmt.Sprintf("%s.%d", object, i), perShard)
		if e != nil {
			return nil, e
		}
	}
	enc := NewEncoder(writers)

	return &RSPutStream{dataServers, object, perShard, writers, enc}, nil
}

func (s *RSPutStream) Write(p []byte) (int, error) {
	return s.enc.Write(p)
}

func (s *RSPutStream) Commit(success bool) {
	s.enc.Close()
	for i := range s.writers {
		s.writers[i].(*objectstream.TempPutStream).Commit(success)
	}
}
