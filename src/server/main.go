package main

import (
	"receive"

	"gopkg.in/gin-gonic/gin.v1"

	"net/http"
	"log"
	"time"
	"path/filepath"
	"os"
	"strings"
)

var wrkDir string

type Server struct {
	srv *http.Server
	engine *gin.Engine
}

func NewServer(addr string) Server {
	router := gin.Default()
	router.LoadHTMLGlob(filepath.Join(wrkDir,"./templates/**/*"))
	router.GET("/*filename", receive.NewService(5*time.Second).Serve)

	srv:= &http.Server{
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 15 * time.Second,

		Addr: addr,
		Handler: router,
	}
	return Server {
		srv:srv,
		engine: router,
	}
}

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	dirs := strings.Split(dir, "/")
	index := len(dirs) - 1
	if dirs[index] == "bin" {
		dirs = append(dirs[:index], dirs[index+1:]...)
		wrkDir = strings.Join(dirs, "/")
	} else {
		wrkDir = dir
	}
}

func main() {
	server := NewServer(":5000")

	log.Fatal(server.srv.ListenAndServe())
}
