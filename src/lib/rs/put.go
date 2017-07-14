package rs

import (
	"fmt"
	"github.com/klauspost/reedsolomon"
	"io"
	"lib/objectstream"
)

const (
	DATA_SHARDS   = 4
	PARITY_SHARDS = 2
	ALL_SHARDS    = DATA_SHARDS + PARITY_SHARDS
)

type RSPutStream struct {
	dataServers   []string
	object        string
	size          int64
	dataWriters   []io.Writer
	parityWriters []io.Writer
	writer        io.Writer
	c             chan error
	enc           reedsolomon.StreamEncoder
}

func NewRSPutStream(dataServers []string, object string, size int64) (*RSPutStream, error) {
	if len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	enc, e := reedsolomon.NewStreamC(DATA_SHARDS, PARITY_SHARDS, true, true)
	if e != nil {
		return nil, e
	}

	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	dataWriters := make([]io.Writer, DATA_SHARDS)
	for i := range dataWriters {
		dataWriters[i], e = objectstream.NewTempPutStream(dataServers[i],
			fmt.Sprintf("%s.%d", object, i), perShard)
		if e != nil {
			return nil, e
		}
	}

	reader, writer := io.Pipe()
	c := make(chan error)
	go func() {
		e := enc.Split(reader, dataWriters, size)
		c <- e
	}()

	parityWriters := make([]io.Writer, PARITY_SHARDS)
	return &RSPutStream{dataServers, object, perShard, dataWriters, parityWriters, writer, c, enc}, nil
}

func (w *RSPutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *RSPutStream) Close() error {
	w.writer.(*io.PipeWriter).Close()
	e := <-w.c
	if e != nil {
		return e
	}

	dataReaders := make([]io.Reader, DATA_SHARDS)
	for i := range dataReaders {
		dataReaders[i], e = w.dataWriters[i].(*objectstream.TempPutStream).NewTempGetStream()
		if e != nil {
			return e
		}
	}
	for i := range w.parityWriters {
		w.parityWriters[i], e = objectstream.NewTempPutStream(w.dataServers[i+DATA_SHARDS],
			fmt.Sprintf("%s.%d", w.object, i+DATA_SHARDS), w.size)
		if e != nil {
			return e
		}
	}

	return w.enc.Encode(dataReaders, w.parityWriters)
}

func (w *RSPutStream) Commit(success bool) {
	w.Close()
	for i := range w.dataWriters {
		if w.dataWriters[i] != nil {
			w.dataWriters[i].(*objectstream.TempPutStream).Commit(success)
		}
	}
	for i := range w.parityWriters {
		if w.parityWriters[i] != nil {
			w.parityWriters[i].(*objectstream.TempPutStream).Commit(success)
		}
	}
}
