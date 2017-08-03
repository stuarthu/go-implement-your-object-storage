package rs

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func testEncodeDecode(t *testing.T, p []byte) {
	writers := make([]io.Writer, ALL_SHARDS)
	readers := make([]io.Reader, ALL_SHARDS)
	for i := range writers {
		writers[i], _ = os.Create(fmt.Sprintf("/tmp/ut_%d", i))
	}
	enc := NewEncoder(writers)
	length := len(p)
	for count := 0; count != length; {
		n, e := enc.Write(p[count:])
		if e != nil {
			t.Error(e)
		}
		count += n
	}
	enc.Close()
	for i := range writers {
		writers[i].(*os.File).Close()
		writers[i] = nil
		readers[i], _ = os.Open(fmt.Sprintf("/tmp/ut_%d", i))
	}
	readers[1] = nil
	readers[4] = nil
	writers[1], _ = os.Create("/tmp/repair_1")
	writers[4], _ = os.Create("/tmp/repair_4")
	dec := NewDecoder(readers, writers, int64(length))
	b := make([]byte, length+10)
	count := 0
	for {
		n, e := dec.Read(b[count:])
		count += n
		if e == io.EOF {
			break
		}
	}
	if count != length {
		t.Error(count, length)
	}
	if !reflect.DeepEqual(b[:count], p) {
		t.Error("not match")
	}
	output, e := exec.Command("diff", "/tmp/ut_1", "/tmp/repair_1").Output()
	if len(output) != 0 {
		t.Error(output, e)
	}
	output, e = exec.Command("diff", "/tmp/ut_4", "/tmp/repair_4").Output()
	if len(output) != 0 {
		t.Error(output, e)
	}
}

func TestEncodeDecode(t *testing.T) {
	p := []byte{1}
	testEncodeDecode(t, p)
	p = []byte("123")
	testEncodeDecode(t, p)
	p = []byte("12345")
	testEncodeDecode(t, p)
	p = make([]byte, 9999)
	fillRandom(p)
	testEncodeDecode(t, p)
	p = make([]byte, 99999)
	fillRandom(p)
	testEncodeDecode(t, p)
}

func fillRandom(p []byte) {
	for i := 0; i < len(p); i += 7 {
		val := rand.Int63()
		for j := 0; i+j < len(p) && j < 7; j++ {
			p[i+j] = byte(val)
			val >>= 8
		}
	}
}
