package tag

import (
	"slices"
	"strings"

	"github.com/EscaperSK/go-notes/lib/app/note"
)

func All(notes note.Notes) []string {
	tags := []string{}
	for _, note := range notes {
		tags = append(tags, note.Tags...)
	}

	slices.Sort(tags)

	return slices.Compact(tags)
}

func Filter(tags []string, search string) []string {
	if len(search) == 0 {
		return tags
	}

	filtered := []string{}

	for _, tag := range tags {
		if !strings.Contains(tag, search) {
			continue
		}

		filtered = append(filtered, tag)
	}

	return filtered
}

func Except(tags []string, except []string) []string {
	result := []string{}

	for _, tag := range tags {
		if slices.Contains(except, tag) {
			continue
		}

		result = append(result, tag)
	}

	return result
}
