package objects

import (
	"../locate"
	"fmt"
	"net/http"
)

func listObjects(w http.ResponseWriter) {
	for object, id := range locate.Objects {
		w.Write([]byte(fmt.Sprintf("%s.%d\n", object, id)))
	}
}
