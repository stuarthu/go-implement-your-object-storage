package objects

import (
	"../heartbeat"
	"../locate"
	"io"
	"lib/es"
	"lib/rs"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	versionId := r.URL.Query()["version"]
	var version es.Metadata
	var e error
	if len(versionId) == 0 {
		version, e = es.SearchLatestVersion(name)
	} else {
		vid, e := strconv.Atoi(versionId[0])
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		version, e = es.GetVersion(name, vid)
	}
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if version.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	object := url.PathEscape(version.Hash)
	locateInfo := locate.Locate(object)
	if len(locateInfo) == 0 {
		log.Println("found object in metadataServer but not in dataServer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ds := make([]string, 0)
	if len(locateInfo) != rs.ALL_SHARDS {
		ds = heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	stream, e := rs.NewRSGetStream(locateInfo, ds, object, version.Size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(w, stream)
}
