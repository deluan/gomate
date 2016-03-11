package gomate

import (
	"fmt"
	"strings"
)

const (
	DefaultNamespace = "gomate-index"
	KeyStoreSuffix   = "all-keys"
)

type Indexer interface {
	Index(id string, doc string) error
	Clear() error
}

type indexer struct {
	namespace string
	keyStore  string
	db        DB
}

func NewIndexer(db DB, namespace ...string) Indexer {
	i := &indexer{db: db, namespace: DefaultNamespace}

	if len(namespace) > 0 {
		i.namespace = namespace[0]
	}
	i.keyStore = fmt.Sprintf("%s:%s", i.namespace, KeyStoreSuffix)

	return i
}

func (i indexer) Index(id string, doc string) error {
	doc = strings.TrimSpace(doc)
	terms := strings.Split(doc, " ")

	for _, t := range terms {
		p := ScorePair{Score: 0, Member: id}
		err := i.db.Zadd(keyForTerm(i.namespace, t), p)
		if err != nil {
			return err
		}
		for _, s := range generatePrefixes(t) {
			p := ScorePair{Score: 1, Member: id}
			err := i.db.Zadd(keyForTerm(i.namespace, s), p)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (i indexer) Clear() error {
	keys, err := i.db.Smembers(i.keyStore)
	if err != nil {
		return err
	}
	for _, k := range keys {
		_, err := i.db.Zclear(k)
		if err != nil {
			return err
		}
	}

	_, err = i.db.Sclear(i.keyStore)
	return err
}

func (i indexer) collectKeys(key string) {
	i.db.Sadd(i.keyStore, key)
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

func keyForTerm(ns string, term string) string {
	return fmt.Sprintf("%s:terms:%s", ns, term)
}

func keyForCache(ns string, term string) string {
	return fmt.Sprintf("%s:cache:%s", ns, term)
}
