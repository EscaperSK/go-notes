package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/EscaperSK/go-notes/lib/app/note"
	"github.com/EscaperSK/go-notes/lib/app/tag"
	"github.com/EscaperSK/go-notes/lib/fs"
)

var templates *template.Template

var notes []note.Note
var tags []string

func Serve() {
	templates = parseTemplates()

	notes = note.All()
	tags = tag.All()

	regHandlers()

	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}

func regHandlers() {
	publicFS := fs.NewPublicFS()

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			data := struct {
				Notes []note.Note
				Tags  []string
			}{notes, tags}
			render(w, "pages.home", data)
			return
		}

		publicFS.ServeHTTP(w, r)
	})
}

func render(w http.ResponseWriter, name string, data any) {
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
