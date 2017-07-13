package temp

import (
	"../heartbeat"
	"encoding/base64"
	"encoding/json"
	"lib/objectstream"
	"lib/utils"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type tempToken struct {
	Name   string
	Size   int64
	Hash   string
	Server string
	Uuid   string
}

func post(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	hash := utils.GetHashFromHeader(r)
	ds := heartbeat.ChooseRandomDataServers(1, nil)
	if len(ds) != 1 {
		log.Println("not enough dataServer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	server := ds[0]
	stream, e := objectstream.NewTempPutStream(server, name, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	uuid := stream.Uuid
	t := tempToken{name, size, hash, server, uuid}
	w.Write([]byte(url.PathEscape(t.toString())))
}

func (t *tempToken) toString() string {
	b, _ := json.Marshal(t)
	return base64.StdEncoding.EncodeToString(b)
}

func tokenFromString(token string) (*tempToken, error) {
	b, e := base64.StdEncoding.DecodeString(token)
	if e != nil {
		return nil, e
	}
	var t tempToken
	e = json.Unmarshal(b, &t)
	return &t, e
}
