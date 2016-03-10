package gomate

import "github.com/siddontang/ledisdb/ledis"

type LedisEmbeddedDB struct {
	db ledis.DB
}

func NewLedisEmbeddedDB(db ledis.DB) DB {
	return &LedisEmbeddedDB{db: db}
}

func (l LedisEmbeddedDB) Zadd(key string, pairs ...ScorePair) error {
	ps := make([]ledis.ScorePair, len(pairs))
	for i, p := range pairs {
		ps[i] = ledis.ScorePair{Score: p.Score, Member: []byte(p.Member)}
	}
	_, err := l.db.ZAdd([]byte(key), ps...)
	return err
}

var _ DB = (*LedisEmbeddedDB)(nil)