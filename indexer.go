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
	keyChain  string
	db        DB
}

const (
	KindSet  = "S-"
	KindZSet = "Z-"
)

func NewIndexer(db DB, namespace ...string) Indexer {
	i := &indexer{db: db, namespace: DefaultNamespace}

	if len(namespace) > 0 {
		i.namespace = namespace[0]
	}
	i.keyChain = fmt.Sprintf("%s:%s", i.namespace, KeyStoreSuffix)

	return i
}

func (i indexer) Index(id string, doc string) error {
	doc = strings.TrimSpace(doc)
	terms := strings.Split(doc, " ")

	for _, t := range terms {
		if err := i.addTerm(id, t, 0); err != nil {
			return err
		}

		for _, s := range generatePrefixes(t) {
			if err := i.addTerm(id, s, 1); err != nil {
				return err
			}
		}
	}
	return nil
}

func (i indexer) addTerm(id string, s string, score int64) error {
	p := ScorePair{Score: score, Member: id}
	sKey := keyForTerm(i.namespace, s)
	err := i.db.Zadd(sKey, p)
	if err != nil {
		return err
	}
	i.collectKeys(sKey, KindZSet)
	return nil
}

func (i indexer) Clear() error {
	keys, err := i.db.Smembers(i.keyChain)
	if err != nil {
		return err
	}
	for _, k := range keys {
		parts := strings.SplitAfterN(k, "-", 2)
		var err error
		switch parts[0] {
		case KindZSet:
			_, err = i.db.Zclear(k)
		case KindSet:
			_, err = i.db.Sclear(k)
		}

		if err != nil {
			return err
		}
	}

	_, err = i.db.Sclear(i.keyChain)
	return err
}

func (i indexer) collectKeys(key string, kind string) {
	i.db.Sadd(i.keyChain, kind+key)
}

func generatePrefixes(term string) []string {
	l := len(term)
	if l < 2 {
		return nil
	}

	ps := make([]string, 0, l-1)
	for i := 1; i < l; i++ {
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
