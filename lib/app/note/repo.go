package note

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
)

const path = "data/storage.json"

func All() []note {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return []note{}
	}

	notes := []note{}
	err = json.Unmarshal(bytes, &notes)
	if err != nil {
		return []note{}
	}

	return notes
}

func Filter(notes []note, filters Filters) []note {
	filtered := []note{}

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
