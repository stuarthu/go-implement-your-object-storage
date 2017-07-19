package objects

import (
	"../locate"
	"os"
	"path/filepath"
)

func del(name string) {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + name + ".*")
	if len(files) != 1 {
		return
	}
	locate.Del(name)
	os.Remove(files[0])
}
