package main

import (
	"receive"

	"gopkg.in/gin-gonic/gin.v1"

	"net/http"
	"log"
	"path/filepath"
)


type Server struct {
	srv *http.Server
	engine *gin.Engine
	settings Settings
}

func NewServer(settings Settings) Server {
	router := gin.Default()
	router.LoadHTMLGlob(filepath.Join(settings.WrkDir,"./templates/**/*"))
	router.GET("/*filename", receive.NewService(settings.Receiever).Serve)

	srv:= &http.Server{
		ReadTimeout: settings.ReadTimeout,
		WriteTimeout: settings.WriteTimeout,

		Addr: settings.ServerAddr,
		Handler: router,
	}
	return Server {
		srv:srv,
		engine: router,
		settings: settings,
	}
}

func main() {
	settings,err := NewSettings()
	if err != nil {
		log.Fatal(err)
	}
	server := NewServer(settings)

	log.Fatal(server.srv.ListenAndServe())
}
