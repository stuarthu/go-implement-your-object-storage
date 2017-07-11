package temp

import (
	"../locate"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	b, e := ioutil.ReadFile(infoFile)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	datFile := infoFile + ".dat"
	f, e := os.OpenFile(datFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	i := strings.Split(string(b), ":")
	size, _ := strconv.ParseInt(i[1], 0, 64)
	actual := info.Size()
	os.Remove(infoFile)
	if actual != size {
		os.Remove(datFile)
		log.Println("actual size mismatch")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
    file := strings.Split(i[0], ".")
    id, _ := strconv.Atoi(file[1])
	locate.Add(file[0], id)
	os.Rename(datFile, os.Getenv("STORAGE_ROOT")+"/objects/"+i[0])
}
