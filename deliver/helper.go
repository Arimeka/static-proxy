package deliver

import (
	"log"
	"os"
	"path/filepath"
)

func isCached(cachedPath string) bool {
	if _, err := os.Stat(cachedPath); err == nil {
		return true
	} else {
		return false
	}
}

func createFile(data *[]byte, dir, filename string) (fullpath string) {
	if dir != "." {
		os.MkdirAll("cache"+string(filepath.Separator)+dir, 0777)
	}

	fullpath = "cache" + string(filepath.Separator) + filename

	fi, err := os.Create(fullpath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if _, err := fi.Write(*data); err != nil {
		log.Fatal(err)
	}
	return
}
