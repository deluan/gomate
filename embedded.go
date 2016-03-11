package gomate

import (
	"github.com/siddontang/ledisdb/ledis"
)

type LedisEmbeddedDB struct {
	db *ledis.DB
}

func NewLedisEmbeddedDB(db *ledis.DB) DB {
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

func (l LedisEmbeddedDB) Zrange(key string, start int, stop int) ([]ScorePair, error) {
	res, err := l.db.ZRange([]byte(key), start, stop)
	if err != nil {
		return nil, err
	}
	ps := make([]ScorePair, len(res))
	for i, p := range res {
		ps[i] = ScorePair{Score: p.Score, Member: string(p.Member)}
	}
	return ps, err
}

var _ DB = (*LedisEmbeddedDB)(nil)
