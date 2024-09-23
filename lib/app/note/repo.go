package note

import (
	"encoding/json"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

const path = "data/storage.json"

func New(name string, content string, tags []string) Note {
	if tags == nil {
		tags = []string{}
	}

	return &note{
		Name:    name,
		Content: content,
		Tags:    tags,
	}
}

func All() Notes {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Notes{}
	}

	values := []note{}
	err = json.Unmarshal(bytes, &values)
	if err != nil {
		return Notes{}
	}

	notes := Notes{}
	for _, note := range values {
		notes = append(notes, &note)
	}

	return notes
}

func Single(noteId int, notes Notes) Note {
	for _, note := range notes {
		if note.Id == noteId {
			return note
		}
	}

	return nil
}

func Filter(notes Notes, filters Filters) Notes {
	filtered := Notes{}

	for _, note := range notes {
		if len(filters.Name) > 0 && !strings.Contains(note.Name, filters.Name) {
			continue
		}
		if len(filters.Tags) > 0 && !containsAll(note.Tags, filters.Tags) {
			continue
		}

		filtered = append(filtered, note)
	}

	return filtered
}

func containsAll(tags []string, elements []string) bool {
	for _, tag := range elements {
		if !slices.Contains(tags, tag) {
			return false
		}
	}

	return true
}

func Create(notes Notes, r *http.Request) Note {
	r.ParseForm()

	postTags := r.PostForm["tags"]
	if postTags == nil {
		postTags = []string{}
	}

	newNote := &note{
		Id:        notes[len(notes)-1].Id + 1,
		Timestamp: time.Now().Unix(),

		Name:    r.PostFormValue("name"),
		Content: r.PostFormValue("content"),
		Tags:    postTags,
	}

	return newNote
}

func Edit(note Note, r *http.Request) {
	r.ParseForm()

	postTags := r.PostForm["tags"]
	if postTags == nil {
		postTags = []string{}
	}

	note.Name = r.PostFormValue("name")
	note.Content = r.PostFormValue("content")
	note.Tags = postTags
}

func Save(notes Notes) error {
	storage, err := os.Create(path)
	if err != nil {
		return err
	}
	defer storage.Close()

	asJson, err := json.MarshalIndent(notes, "", "\t")
	if err != nil {
		return err
	}

	_, err = storage.Write(asJson)
	if err != nil {
		return err
	}

	return nil
}

func Delete(noteId int, notes Notes) Notes {
	return slices.DeleteFunc(notes, func(n Note) bool {
		return n.Id == noteId
	})
}
