package server

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/EscaperSK/go-notes/lib/app/note"
	"github.com/EscaperSK/go-notes/lib/app/tag"
	"github.com/EscaperSK/go-notes/lib/fs"
)

var layouts *template.Template
var templates *template.Template

var notes note.Notes
var tags []string

func Serve() {
	layouts = parseLayouts()
	templates = parseTemplates()

	notes = note.All()
	tags = tag.All()

	regHandlers()

	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}

func regHandlers() {
	http.Handle("GET /", fs.NewPublicFS())

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Notes note.Notes
			Tags  []string
		}{notes, tags}

		render(w, "layout", "pages.home", data)
		return
	})

	http.HandleFunc("GET /filter", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		search := query.Get("search")
		tags := query["tags"]

		filters := note.Filters{Name: search, Tags: tags}
		data := note.Filter(notes, filters)

		render(w, "layout", "note.list", data)
	})

	http.HandleFunc("GET /note/{noteId}", func(w http.ResponseWriter, r *http.Request) {
		pathNoteId := r.PathValue("noteId")

		noteId, err := strconv.Atoi(pathNoteId)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		single := note.Single(noteId, notes)
		if single == nil {
			http.NotFound(w, r)
			return
		}

		data := struct {
			Note note.Note
			Tags []string
		}{single, tags}

		render(w, "layout", "pages.note", data)
	})
}

func render(w http.ResponseWriter, layout string, name string, data any) {
	buf := bytes.NewBuffer([]byte{})
	err := templates.ExecuteTemplate(buf, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	content := template.HTML(buf.Bytes())

	err = layouts.ExecuteTemplate(w, layout, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
