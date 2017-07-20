package objects

import (
	"../locate"
	"os"
	"path/filepath"
	"strings"
)

func del(name string) {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + name + ".*")
	if len(files) != 1 {
		return
	}
	object := strings.Split(name, ".")[0]
	locate.Del(object)
	os.Rename(files[0], os.Getenv("STORAGE_ROOT")+"/garbage/"+filepath.Base(files[0]))
}
