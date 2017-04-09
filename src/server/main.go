package main

import (
	"receive"

	"gopkg.in/gin-gonic/gin.v1"

	"net/http"
	"log"
	"time"
)

type Server struct {
	srv *http.Server
	engine *gin.Engine
}

func NewServer(addr string) Server {
	router := gin.Default()
	router.GET("/", receive.NewService(5*time.Second).Serve)

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

func main() {
	server := NewServer(":5000")

	log.Fatal(server.srv.ListenAndServe())
}
