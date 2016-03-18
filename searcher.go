package gomate

import (
	"sort"
	"strings"
)

const CacheTimeOut = 600

type Searcher interface {
	Search(query string, min int, max int) ([]string, error)
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

func (s searcher) Search(query string, min int, max int) ([]string, error) {
	query = strings.TrimSpace(query)
	var resp []ScorePair
	var err error
	var finalIdx string

	terms := strings.Split(query, " ")

	finalIdx, err = s.multiWordQuery(terms)
	if err != nil || finalIdx == "" {
		return nil, err
	}

	resp, err = s.db.Zrange(finalIdx, min, max)
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
	sort.Strings(terms)
	final := strings.Join(terms, "|")
	finalIdx := keyForCache(s.namespace, final)

	exists, err := s.db.Zkeyexists(finalIdx)
	if exists && err == nil {
		return finalIdx, nil
	}

	idxs := make([]string, len(terms)+1)
	idxs[0] = idSetName(s.namespace)
	for i, t := range terms {
		idxs[i+1] = keyForTerm(s.namespace, t)
	}
	r, err := s.db.Zinterstore(finalIdx, idxs, AggregateSum)
	if err != nil {
		return "", err
	}
	if r == 0 {
		return "", nil
	}
	s.db.Zexpire(finalIdx, CacheTimeOut)
	return finalIdx, nil
}
