package gomate

import "errors"

type Searcher interface {
	Search(query string) ([]string, error)
}

type searcher struct {
	namespace string
	db        DB
}

func NewSearcher(db DB, namespace ...string) Searcher {
	i := &searcher{db: db, namespace: defaultNamespace}

	if len(namespace) > 0 {
		i.namespace = namespace[0]
	}

	return i
}

func (s searcher) Search(query string) ([]string, error) {
	return nil, errors.New("not implemented")
}
