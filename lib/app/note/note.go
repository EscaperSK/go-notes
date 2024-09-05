package note

const gistLen = 140

type Note struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (n *Note) Gist() string {
	if len(n.Content) <= gistLen+3 {
		return n.Content
	}

	return n.Content[:gistLen] + "..."
}
