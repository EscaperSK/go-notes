package note

import "time"

const gistLen = 140

type Note *note
type Notes []Note

type note struct {
	Id        int   `json:"id"`
	Timestamp int64 `json:"timestamp"`

	Name    string   `json:"name"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (n *note) Gist() string {
	if len(n.Content) <= gistLen+3 {
		return n.Content
	}

	return n.Content[:gistLen] + "..."
}

func (n *note) DateString() string {
	if n.Timestamp == 0 {
		return "Дата не указана"
	}

	return time.Unix(n.Timestamp, 0).Format("15:04 02.01.06")
}

func (n *note) FullDateString() string {
	if n.Timestamp == 0 {
		return "Дата не указана"
	}

	return time.Unix(n.Timestamp, 0).Format("15:04:05 02.01.2006")
}
