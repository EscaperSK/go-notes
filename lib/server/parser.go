package server

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

const (
	PathSep string = string(os.PathSeparator)
	PartSep string = "."
)

func parseTemplates() *template.Template {
	templates := template.New("")
	root := filepath.Clean("lib/templates")
	skip := len(root) + 1

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}
		if err != nil {
			return err
		}

		bytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		t := templates.New(getName(path, skip))

		_, err = t.Parse(string(bytes))
		if err != nil {
			return err
		}

		return nil
	})

	return template.Must(templates, err)
}

func getName(path string, skip int) string {
	file := path[skip:]
	ext := filepath.Ext(file)
	clean := strings.TrimSuffix(file, ext)
	parts := strings.Split(clean, PathSep)

	return strings.Join(parts, PartSep)
}
