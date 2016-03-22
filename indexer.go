package gomate

import (
	"fmt"
	"strings"
)

const (
	DefaultNamespace = "gomate-index"
	IdSetSuffix      = "all-ids"
	KeyChainSuffix   = "all-keys"
)

type Indexer interface {
	Index(id string, doc string) error
	Remove(ids ...string) error
	Clear() error
}

type indexer struct {
	namespace string
	keyChain  string
	idSet     string
	db        DB
}

const (
	KindSet  = "S-"
	KindZSet = "Z-"
	KindHash = "H-"
)

func NewIndexer(db DB, namespace ...string) Indexer {
	i := &indexer{db: db, namespace: DefaultNamespace}

	if len(namespace) > 0 {
		i.namespace = namespace[0]
	}
	i.keyChain = keyChainName(i.namespace)
	i.idSet = idSetName(i.namespace)

	return i
}

func keyChainName(namespace string) string {
	return fmt.Sprintf("%s:%s", namespace, KeyChainSuffix)
}

func idSetName(namespace string) string {
	return fmt.Sprintf("%s:%s", namespace, IdSetSuffix)
}

func (i *indexer) Index(id string, doc string) error {
	doc = strings.TrimSpace(doc)

	// First remove the id if it was already indexed
	i.Remove(id)

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

		if err := i.addId(id, doc); err != nil {
			return err
		}
	}
	return nil
}

func (i *indexer) addId(id string, doc string) error {
	_, err := i.db.Hset(i.idSet, id, doc)
	return err
}

func (i *indexer) addTerm(id string, s string, score int64) error {
	p := ScorePair{Score: score, Member: id}
	sKey := keyForTerm(i.namespace, s)
	err := i.db.Zadd(sKey, p)
	if err != nil {
		return err
	}
	i.collectKeys(sKey, KindZSet)
	return nil
}

func (i *indexer) Remove(ids ...string) error {
	docs, err := i.db.Hmget(i.idSet, ids...)
	if err != nil {
		return err
	}
	temp := make(map[string][]string)

	for i, doc := range docs {
		if doc == "" {
			continue
		}
		terms := strings.Split(doc, " ")
		id := ids[i]
		for _, t := range terms {
			temp[t] = append(temp[t], id)
			prefixes := generatePrefixes(t)
			for _, s := range prefixes {
				temp[s] = append(temp[s], id)
			}
		}
	}

	for s, l := range temp {
		i.db.Zrem(keyForTerm(i.namespace, s), l...)
	}
	return nil
}

func (i *indexer) Clear() error {
	keys, err := i.db.Smembers(i.keyChain)
	if err != nil {
		return err
	}
	for _, k := range keys {
		parts := strings.SplitAfterN(k, "-", 2)
		var err error
		switch parts[0] {
		case KindZSet:
			_, err = i.db.Zclear(parts[1])
		case KindSet:
			_, err = i.db.Sclear(parts[1])
		case KindHash:
			_, err = i.db.Hclear(parts[1])
		}

		if err != nil {
			return err
		}
	}

	_, err = i.db.Hclear(i.idSet)
	_, err = i.db.Sclear(i.keyChain)
	return err
}

func (i *indexer) collectKeys(key string, kind string) {
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
