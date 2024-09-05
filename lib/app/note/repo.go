package note

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
)

const path = "data/storage.json"

func All() []Note {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return []Note{}
	}

	notes := []Note{}
	err = json.Unmarshal(bytes, &notes)
	if err != nil {
		return []Note{}
	}

	return notes
}

func Filter(notes []Note, filters Filters) []Note {
	filtered := []Note{}

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
