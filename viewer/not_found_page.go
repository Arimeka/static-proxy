package viewer

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

func NotFoundPage() http.HandlerFunc {

	path, err := filepath.Abs("./templates/404.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl, templateError := template.ParseFiles(path)
	if templateError != nil {
		log.Fatal(templateError)
	}

	return func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(404)
		tmpl.Execute(response, nil)
	}
}
