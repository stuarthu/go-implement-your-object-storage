package temp

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func post(w http.ResponseWriter, r *http.Request) {
	output, _ := exec.Command("uuidgen").Output()
	uuid := strings.TrimSuffix(string(output), "\n")
	f, e := os.Create(os.Getenv("TMP_ROOT") + "/" + uuid)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	filename := strings.Split(r.URL.EscapedPath(), "/")[2]
	size := r.Header.Get("size")
	fmt.Fprintf(f, "%s:%s", filename, size)
	w.Write([]byte(uuid))
}
