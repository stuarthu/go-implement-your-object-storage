package rs

import (
	"fmt"
	"io"
	"lib/objectstream"
)

type RSGetStream struct {
	writers []io.Writer
	dec     *decoder
}

func NewRSGetStream(locateInfo map[int]string, dataServers []string, object string, size int64) (*RSGetStream, error) {
	if len(locateInfo)+len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	readers := make([]io.Reader, ALL_SHARDS)
	writers := make([]io.Writer, ALL_SHARDS)
	totalReaders := ALL_SHARDS
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	var e error
	for i := 0; i < ALL_SHARDS; i++ {
		server := locateInfo[i]
		if server == "" {
			readers[i] = nil
			totalReaders--
			server = dataServers[0]
			dataServers = dataServers[1:]
			locateInfo[i] = server
			writers[i], e = objectstream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", object, i), perShard)
			if e != nil {
				return nil, e
			}
		} else {
			readers[i], e = objectstream.NewGetStream(server, fmt.Sprintf("%s.%d", object, i))
			if e != nil {
				readers[i] = nil
				totalReaders--
				writers[i], e = objectstream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", object, i), perShard)
				if e != nil {
					return nil, e
				}
			}
		}
	}

	dec := NewDecoder(readers, writers, size)
	return &RSGetStream{writers, dec}, nil
}

func (s *RSGetStream) Read(p []byte) (n int, err error) {
	return s.dec.Read(p)
}

func (s *RSGetStream) Close() {
	for i := range s.writers {
		if s.writers[i] != nil {
			s.writers[i].(*objectstream.TempPutStream).Commit(true)
		}
	}
}
