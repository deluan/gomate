package gomate

import (
	"fmt"
	"strings"
)

const defaultNamespace = "gomate-index"

type Indexer interface {
	Index(id string, doc string) error
}

type indexer struct {
	namespace string
	db        DB
}

func NewIndexer(db DB, namespace ...string) Indexer {
	i := &indexer{db: db, namespace: defaultNamespace}

	if len(namespace) > 0 {
		i.namespace = namespace[0]
	}

	return i
}

func (i indexer) Index(id string, doc string) error {
	terms := strings.Split(doc, " ")

	for _, t := range terms {
		p := ScorePair{Score: 1, Member: id}
		err := i.db.Zadd(i.keyForTerm(t), p)
		if err != nil {
			return err
		}
		for _, s := range generatePrefixes(t) {
			p := ScorePair{Score: 0, Member: id}
			err := i.db.Zadd(i.keyForTerm(s), p)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func generatePrefixes(term string) []string {
	l := len(term)
	if l <= 2 {
		return []string{}
	}

	ps := make([]string, 0, l-1)
	for i := 2; i < l; i++ {
		ps = append(ps, term[0:i])
	}

	return ps
}

func (i indexer) keyForTerm(term string) string {
	return fmt.Sprintf("%s:terms:%s", i.namespace, term)
}
