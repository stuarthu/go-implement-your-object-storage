package rs

import (
	"fmt"
	"github.com/klauspost/reedsolomon"
	"io"
	"net/http"
	"sync"
)

const (
	DATA_SHARDS   = 4
	PARITY_SHARDS = 2
	ALL_SHARDS    = DATA_SHARDS + PARITY_SHARDS
)

func Put(data io.Reader, dataServers []string, object string, size int64) error {
	enc, e := reedsolomon.NewStreamC(DATA_SHARDS, PARITY_SHARDS, true, true)
	if e != nil {
		return e
	}

	readers := make([]io.Reader, DATA_SHARDS)
	writers := make([]io.Writer, DATA_SHARDS)
	for i := range readers {
		readers[i], writers[i] = io.Pipe()
	}
	fmt.Println("split")
	var wg sync.WaitGroup
	wg.Add(DATA_SHARDS)
	go func() {
		enc.Split(data, writers, size)
		for i := range writers {
			writers[i].(*io.PipeWriter).Close()
		}
	}()
	for index := range readers {
		go func(i int) {
			defer wg.Done()
            url := fmt.Sprintf("http://%s/objects/%s:%d", dataServers[i], object, i)
			request, e := http.NewRequest("PUT", url, readers[i])
			if e != nil {
				fmt.Println(e)
				return
			}
			client := http.Client{}
			r, e := client.Do(request)
			fmt.Println(r.StatusCode)
			if e != nil {
				fmt.Println(e)
				return
			}
			r, e = client.Get(url)
			if e != nil {
				fmt.Println(e)
				return
			}
			readers[i] = r.Body
		}(index)
	}
	wg.Wait()
	if e != nil {
		return e
	}

	parityWriters := make([]io.Writer, PARITY_SHARDS)
	parityReaders := make([]io.Reader, PARITY_SHARDS)
	for i := range parityReaders {
		parityReaders[i], parityWriters[i] = io.Pipe()
	}
	fmt.Println("encode")
	wg.Add(PARITY_SHARDS)
	go func() {
		enc.Encode(readers, parityWriters)
		for i := range parityWriters {
			parityWriters[i].(*io.PipeWriter).Close()
		}
	}()
	for index := range parityReaders {
		go func(i int) {
			defer wg.Done()
            url := fmt.Sprintf("http://%s/objects/%s:%d", dataServers[DATA_SHARDS+i], object, DATA_SHARDS+i)
			request, e := http.NewRequest("PUT", url, parityReaders[i])
			if e != nil {
				return
			}
			client := http.Client{}
			_, e = client.Do(request)
			if e != nil {
				return
			}
		}(index)
	}
	wg.Wait()
	return e
}
