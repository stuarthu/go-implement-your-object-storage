package objects

import (
	"log"
	"net/http"
)

func put(w http.ResponseWriter, r *http.Request) {
	hash := getHashFromHeader(r)
	if hash == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	e := storeObject(r, hash)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	e = addVersion(r, hash)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func storeObject(r *http.Request, hash string) error {
	s := heartbeat.ChooseRandomDataServer()
	if s == "" {
		return fmt.Errorf("cannot find any dataServer")
	}
	stream := httpstream.NewPutStream("http://" + s + "/objects/" + hash)
	io.Copy(stream, r.Body)
	return stream.Close()
}
