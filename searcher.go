package gomate

import (
	"sort"
	"strings"
)

type Searcher interface {
	Search(query string) ([]string, error)
}

type searcher struct {
	namespace string
	db        DB
}

func NewSearcher(db DB, namespace ...string) Searcher {
	i := &searcher{db: db, namespace: DefaultNamespace}

	if len(namespace) > 0 {
		i.namespace = namespace[0]
	}

	return i
}

func (s searcher) Search(query string) ([]string, error) {
	var resp []ScorePair
	var err error
	var finalIdx string

	terms := strings.Split(query, " ")
	idxs := make([]string, len(terms))
	for i, t := range terms {
		idxs[i] = keyForTerm(s.namespace, t)
	}

	if len(terms) == 1 {
		finalIdx = idxs[0]
	} else {
		sort.Strings(terms)
		final := strings.Join(terms, "|")
		finalIdx = keyForTerm(s.namespace, final)

		r, err := s.db.Zinterstore(finalIdx, idxs, AggregateSum)
		if err != nil {
			return nil, err
		}

		if r == 0 {
			return nil, nil
		}
	}

	resp, err = s.db.Zrevrange(finalIdx, 0, -1)
	if err != nil {
		return nil, err
	}

	r := make([]string, len(resp))
	for i, p := range resp {
		r[i] = p.Member
	}
	return r, nil
}
