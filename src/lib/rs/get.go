package rs

import (
	"fmt"
	"github.com/klauspost/reedsolomon"
	"io"
	"lib/objectstream"
)

type RSGetStream struct {
	reader io.Reader
}

func NewRSGetStream(locateInfo map[int]string, dataServers []string, object string, size int64) (*RSGetStream, error) {
	if len(locateInfo)+len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	dataReaders := make([]io.Reader, ALL_SHARDS)
	totalReaders := ALL_SHARDS
	var e error
	for i := 0; i < ALL_SHARDS; i++ {
		server := locateInfo[i]
		if server == "" {
			dataReaders[i] = nil
			totalReaders--
			server = dataServers[0]
			dataServers = dataServers[1:]
			locateInfo[i] = server
		} else {
			dataReaders[i], e = objectstream.NewGetStream(server, fmt.Sprintf("%s.%d", object, i))
			if e != nil {
				dataReaders[i] = nil
				totalReaders--
			}
		}
	}
	if totalReaders < DATA_SHARDS {
		return nil, fmt.Errorf("not enough shards")
	}
	if totalReaders != ALL_SHARDS {
		e := repair(locateInfo, dataReaders, object, size)
		if e != nil {
			return nil, e
		}
	}

	for i := range dataReaders {
		dataReaders[i], e = objectstream.NewGetStream(locateInfo[i],
			fmt.Sprintf("%s.%d", object, i))
		if e != nil {
			return nil, e
		}
	}

	enc, e := reedsolomon.NewStreamC(DATA_SHARDS, PARITY_SHARDS, true, true)
	if e != nil {
		return nil, e
	}

	reader, writer := io.Pipe()
	go func() {
		enc.Join(writer, dataReaders, size)
		writer.Close()
	}()

	return &RSGetStream{reader}, nil
}

func (w *RSGetStream) Read(p []byte) (n int, err error) {
	return w.reader.Read(p)
}

func repair(locateInfo map[int]string, dataReaders []io.Reader, object string, size int64) error {
	enc, e := reedsolomon.NewStreamC(DATA_SHARDS, PARITY_SHARDS, true, true)
	if e != nil {
		return e
	}

	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	dataWriters := make([]io.Writer, ALL_SHARDS)
	for i := 0; i < ALL_SHARDS; i++ {
		if dataReaders[i] != nil {
			dataWriters[i] = nil
		} else {
			dataWriters[i], e = objectstream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", object, i), perShard)
			if e != nil {
				return e
			}
		}
	}

	e = enc.Reconstruct(dataReaders, dataWriters)
	if e != nil {
		for i := range dataWriters {
			if dataWriters[i] != nil {
				dataWriters[i].(*objectstream.TempPutStream).Commit(false)
			}
		}
		return e
	}
	for i := range dataWriters {
		if dataWriters[i] != nil {
			dataWriters[i].(*objectstream.TempPutStream).Commit(true)
		}
	}
	return nil
}
