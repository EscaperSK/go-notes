package server

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/EscaperSK/go-notes/lib/app/note"
	"github.com/EscaperSK/go-notes/lib/app/tag"
	"github.com/EscaperSK/go-notes/lib/errs"
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
	tags = tag.All(notes)

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

		renderPage(w, "layout", "pages.home", data)
		return
	})

	http.HandleFunc("GET /filter", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		search := query.Get("search")
		tags := query["tags"]

		filters := note.Filters{Name: search, Tags: tags}
		data := note.Filter(notes, filters)

		renderTmpl(w, "note.list", data)
	})

	http.HandleFunc("GET /note/{noteId}", func(w http.ResponseWriter, r *http.Request) {
		pathNoteId := r.PathValue("noteId")

		noteId, err := strconv.Atoi(pathNoteId)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		data := note.Single(noteId, notes)
		if data == nil {
			http.NotFound(w, r)
			return
		}

		renderPage(w, "layout", "note.view", data)
	})

	http.HandleFunc("GET /note/{noteId}/view", func(w http.ResponseWriter, r *http.Request) {
		pathNoteId := r.PathValue("noteId")

		noteId, err := strconv.Atoi(pathNoteId)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		data := note.Single(noteId, notes)
		if data == nil {
			http.NotFound(w, r)
			return
		}

		renderTmpl(w, "note.view", data)
	})

	http.HandleFunc("GET /note", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Note note.Note
			Tags []string
			Errs errs.ValidationErrors
		}{nil, tags, nil}

		renderPage(w, "layout", "note.create", data)
	})

	http.HandleFunc("POST /note", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		name := r.PostFormValue("name")
		content := r.PostFormValue("content")

		validationErrs := errs.ValidationErrors{}
		if len(name) <= 0 {
			validationErrs["name"] = "Это обязательное поле"
		}
		if len(content) <= 0 {
			validationErrs["content"] = "Это обязательное поле"
		}

		if len(validationErrs) > 0 {
			noteTags := r.PostForm["tags"]
			single := note.New(name, content, noteTags)
			otherTags := tag.Except(tags, single.Tags)

			data := struct {
				Note note.Note
				Tags []string
				Errs errs.ValidationErrors
			}{single, otherTags, validationErrs}

			renderPage(w, "layout", "note.create", data)
			return
		}

		newNote := note.Create(notes, r)

		notes = append(notes, newNote)
		tags = tag.All(notes)

		err := note.Save(notes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		url := fmt.Sprintf("/note/%d", newNote.Id)

		http.Redirect(w, r, url, http.StatusFound)
	})

	http.HandleFunc("GET /note/{noteId}/edit", func(w http.ResponseWriter, r *http.Request) {
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

		otherTags := tag.Except(tags, single.Tags)

		data := struct {
			Note note.Note
			Tags []string
			Errs errs.ValidationErrors
		}{single, otherTags, nil}

		renderTmpl(w, "note.edit", data)
	})

	http.HandleFunc("PUT /note/{noteId}", func(w http.ResponseWriter, r *http.Request) {
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

		r.ParseForm()

		name := r.PostFormValue("name")
		content := r.PostFormValue("content")

		validationErrs := errs.ValidationErrors{}
		if len(name) <= 0 {
			validationErrs["name"] = "Это обязательное поле"
		}
		if len(content) <= 0 {
			validationErrs["content"] = "Это обязательное поле"
		}

		if len(validationErrs) > 0 {
			id := single.Id
			timestamp := single.Timestamp
			noteTags := r.PostForm["tags"]

			single := note.New(name, content, noteTags)
			single.Id = id
			single.Timestamp = timestamp

			otherTags := tag.Except(tags, single.Tags)

			data := struct {
				Note note.Note
				Tags []string
				Errs errs.ValidationErrors
			}{single, otherTags, validationErrs}

			renderTmpl(w, "note.edit", data)
			return
		}

		note.Edit(single, r)

		tags = tag.All(notes)

		err = note.Save(notes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		renderTmpl(w, "note.view", single)
	})
}

func renderTmpl(w http.ResponseWriter, name string, data any) {
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderPage(w http.ResponseWriter, layout string, name string, data any) {
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
