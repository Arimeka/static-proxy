package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"runtime"
	"static-proxy/deliver"
	"static-proxy/settings"
	"static-proxy/viewer"
)

const (
	VERSION = "0.0.1"
)

func init() {
	settings.SetupOptions("config address env")
}

func main() {
	settings.Setup()

	config, err := settings.SetDefaults(settings.BuildMain(settings.Path))
	if err != nil {
		log.Fatal(err)
	}

	s3Config, err := settings.BuildS3(settings.S3Path)
	if err != nil {
		log.Fatal(err)
	}

	runtime.GOMAXPROCS(config.NumCPU)

	printBanner(config)

	fs := http.FileServer(http.Dir("./static"))

	// Mandatory root-based resources
	serveSingle("/favicon.ico", "./static/favicon.ico")
	serveSingle("/robots.txt", "./static/robots.txt")

	router := mux.NewRouter()

	router.PathPrefix("/static/").Handler(viewer.ServeStatic(http.StripPrefix("/static/", fs))).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(viewer.NotFoundPage())

	server := deliver.New(config.Workers, s3Config)
	router.HandleFunc("/{filename:.+}", server.Handle()).Methods("GET")

	http.Handle("/", handlers.LoggingHandler(os.Stdout, router))

	log.Fatal(http.ListenAndServe(config.Address+":"+config.Port, nil))
}

func printBanner(config settings.Settings) {
	log.Println("StaticProxy", VERSION, "("+runtime.Version()+" "+runtime.GOOS+"/"+runtime.GOARCH+")")
	log.Println("- environment:", settings.Env)
	log.Println("- numcpu:     ", config.NumCPU)
	log.Println("listen", config.Address+":"+config.Port)
}

func serveSingle(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}
