package viewer

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
)

func ServeStatic(h http.Handler) http.HandlerFunc {

	path, err := filepath.Abs("./templates/404.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl, templateError := template.ParseFiles(path)
	if templateError != nil {
		log.Fatal(templateError)
	}

	return func(response http.ResponseWriter, request *http.Request) {
		if strings.HasSuffix(request.URL.Path, "/") {
			response.WriteHeader(404)
			tmpl.Execute(response, nil)
		} else {
			h.ServeHTTP(response, request)
		}
	}
}
