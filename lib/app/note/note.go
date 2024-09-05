package note

type note struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func New(id int, name string, content string, tags []string) *note {
	return &note{
		Id:      id,
		Name:    name,
		Content: content,
		Tags:    tags,
	}
}
