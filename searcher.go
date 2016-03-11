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
	query = strings.TrimSpace(query)
	var resp []ScorePair
	var err error
	var finalIdx string

	terms := strings.Split(query, " ")

	if len(terms) == 1 {
		finalIdx = keyForTerm(s.namespace, query)
	} else {
		finalIdx, err = s.multiWordQuery(terms)
		if finalIdx == "" {
			return nil, err
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

func (s searcher) multiWordQuery(terms []string) (string, error) {
	idxs := make([]string, len(terms))
	for i, t := range terms {
		idxs[i] = keyForTerm(s.namespace, t)
	}

	sort.Strings(terms)
	final := strings.Join(terms, "|")
	finalIdx := keyForTerm(s.namespace, final)

	r, err := s.db.Zinterstore(finalIdx, idxs, AggregateSum)
	if err != nil {
		return "", err
	}
	if r == 0 {
		return "", nil
	}
	return finalIdx, nil
}
