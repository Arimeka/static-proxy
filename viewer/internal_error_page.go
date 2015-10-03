package viewer

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

func InternalErrorPage(next http.Handler) http.HandlerFunc {

	path, err := filepath.Abs("./templates/500.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl, templateError := template.ParseFiles(path)
	if templateError != nil {
		log.Fatal(templateError)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				w.WriteHeader(500)
				tmpl.Execute(w, nil)
			}
		}()

		next.ServeHTTP(w, r)
	}
}
