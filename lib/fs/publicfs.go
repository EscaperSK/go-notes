package fs

import (
	iofs "io/fs"
	"net/http"
)

type fs struct {
	fs http.FileSystem
}

func (p fs) Open(path string) (http.File, error) {
	f, err := p.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		return nil, iofs.ErrNotExist
	}

	return f, nil
}

func NewPublicFS() http.Handler {
	return http.FileServer(fs{http.Dir("public")})
}
