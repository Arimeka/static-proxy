package main

import (
	"recive"

	"github.com/julienschmidt/httprouter"
	"github.com/Sirupsen/logrus"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"github.com/go-errors/errors"

	"net/http"
	"log"
	"time"
	"os"
)

type Server struct {
	srv *http.Server
	router *httprouter.Router

	Log *logrus.Logger
}

func NewServer(logger *logrus.Logger, addr string) Server {


	router := httprouter.New()


	router.Handler("GET","/",recive.NewHandler(logger, 5*time.Second))

	chain := alice.New(nosurf.NewPure).Then(router)

	srv := &http.Server{
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 15 * time.Second,

		Addr: addr,
		Handler: chain,
	}

	server := Server{
		srv: srv,
		router: router,
		Log: logger,
	}
	router.PanicHandler = server.recover

	return server
}

func main() {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.Formatter = &logrus.TextFormatter{}

	server := NewServer(logger, ":5000")

	log.Fatal(server.srv.ListenAndServe())
}

func (server Server) recover(w http.ResponseWriter, r *http.Request, err interface{}) {
	server.Log.Panic(errors.Wrap(err, 2).ErrorStack())
}
