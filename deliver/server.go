package deliver

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path/filepath"
	"static-proxy/settings"
	"text/template"
)

type Job struct {
	Filename     string
	Result       string
	BucketConfig interface{}
}

type Request struct {
	Job        *Job
	ResultChan chan string
}

type Server struct {
	Requests chan *Request
	S3Config settings.S3
}

func New(pool int, s3Config settings.S3) Server {
	jobs, results := WorkerPool(pool)
	jobs, results = Cache(jobs, results)
	requests := RequestMux(jobs, results)

	return Server{Requests: requests, S3Config: s3Config}
}

func (s Server) Handle() http.HandlerFunc {
	var name string

	path, err := filepath.Abs("./templates/404.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl, templateError := template.ParseFiles(path)
	if templateError != nil {
		log.Fatal(templateError)
	}

	return func(w http.ResponseWriter, req *http.Request) {
		var path string
		vars := mux.Vars(req)
		filename := vars["filename"]

		name = req.Host

		if s.S3Config.Host[name] != nil {
			request := &Request{Job: &Job{Filename: filename, BucketConfig: s.S3Config.Host[name]}, ResultChan: make(chan string)}
			s.Requests <- request
			path = <-request.ResultChan
		}
		if path != "" {
			http.ServeFile(w, req, path)
		} else {
			w.WriteHeader(404)
			tmpl.Execute(w, nil)
		}
	}
}

func RequestMux(jobs chan *Job, results chan *Job) chan *Request {
	requests := make(chan *Request)

	go func() {
		queues := make(map[string][]*Request)

		for {
			select {
			case request := <-requests:
				job := request.Job
				queues[job.Filename] = append(queues[job.Filename], request)

				if len(queues[job.Filename]) == 1 { // the one we appended is the first one
					go func() {
						jobs <- job
					}()
				}

			case job := <-results:
				for _, request := range queues[job.Filename] {
					request.ResultChan <- job.Result
				}

				delete(queues, job.Filename)
			}
		}
	}()

	return requests
}

func Cache(upstreamJobs chan *Job, upstreamResults chan *Job) (chan *Job, chan *Job) {
	jobs := make(chan *Job)
	results := make(chan *Job)

	go func() {
		for {
			select {
			case job := <-jobs:
				cachedPath := "cache" + string(filepath.Separator) + job.Filename
				if isCached(cachedPath) { // cache hit
					log.Printf("Cache hit: %s", job.Filename)
					job.Result = cachedPath
					results <- job
				} else { // cache miss
					log.Printf("Cache miss: %s", job.Filename)
					upstreamJobs <- job
				}

			case job := <-upstreamResults:
				job.Result = job.Result
				results <- job
			}
		}
	}()

	return jobs, results
}
