package rs

import (
	"github.com/klauspost/reedsolomon"
	"io"
)

type encoder struct {
	writers   []io.Writer
	enc       reedsolomon.Encoder
	cache     []byte
	cacheSize int
}

func NewEncoder(writers []io.Writer) *encoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &encoder{writers, enc, nil, 0}
}

func (e *encoder) Write(p []byte) (n int, err error) {
	length := len(p)
	if e.cacheSize+length < BLOCK_SIZE {
		e.cache = append(e.cache, p...)
		e.cacheSize += length
		return length, nil
	}
	length = BLOCK_SIZE - e.cacheSize
	e.cache = append(e.cache, p[:length]...)
	e.flush()
	return length, nil
}

func (e *encoder) flush() {
	shards, _ := e.enc.Split(e.cache)
	e.enc.Encode(shards)
	for i := range shards {
		e.writers[i].Write(shards[i])
	}
	e.cache = []byte{}
	e.cacheSize = 0
}
