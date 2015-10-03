package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"runtime"
	"static-proxy/deliver"
	"static-proxy/settings"
	"static-proxy/utils"
	"static-proxy/viewer"
	"time"
)

const (
	VERSION = "0.0.2"
)

func init() {
	settings.SetupOptions("config address env")
}

func main() {
	settings.Setup()

	err := settings.Build()
	if err != nil {
		log.Fatal(fmt.Errorf("Unable to load configuration, %v", err))
	}

	runtime.GOMAXPROCS(settings.Config.NumCPU)

	printBanner(settings.Config)

	fs := http.FileServer(http.Dir("./static"))

	// Mandatory root-based resources
	serveSingle("/favicon.ico", "./static/favicon.ico")
	serveSingle("/robots.txt", "./static/robots.txt")

	router := mux.NewRouter()
	commonHandlers := alice.New(loggingHandler, recoverHandler)

	router.PathPrefix("/static/").Handler(viewer.ServeStatic(http.StripPrefix("/static/", fs))).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(viewer.NotFoundPage())

	server := deliver.New(settings.Config.Workers)

	router.Handle("/{filename:.+}", server).Methods("GET")

	log.Fatal(http.ListenAndServe(settings.Config.Address+":"+settings.Config.Port, commonHandlers.Then(router)))
}

func printBanner(config *settings.Settings) {
	log.Println("StaticProxy", VERSION, "("+runtime.Version()+" "+runtime.GOOS+"/"+runtime.GOARCH+")")
	log.Println("- environment:", settings.AppSettings.GetString("env"))
	log.Println("- numcpu:     ", config.NumCPU)
	log.Println("listen", config.Address+":"+config.Port)
}

func serveSingle(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("%s - [%s] %q %v\n", utils.GetIpFromRequest(r), r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func recoverHandler(next http.Handler) http.Handler {
	return viewer.InternalErrorPage(next)
}
