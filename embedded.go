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

func (l LedisEmbeddedDB) Zrevrange(key string, start int, stop int) ([]ScorePair, error) {
	res, err := l.db.ZRevRange([]byte(key), start, stop)
	if err != nil {
		return nil, err
	}
	ps := make([]ScorePair, len(res))
	for i, p := range res {
		ps[i] = ScorePair{Score: p.Score, Member: string(p.Member)}
	}
	return ps, nil
}

func (l LedisEmbeddedDB) Zinterstore(destKey string, srcKeys []string, aggregate byte) (int64, error) {
	sk := make([][]byte, len(srcKeys))
	for i, k := range srcKeys {
		sk[i] = []byte(k)
	}
	return l.db.ZInterStore([]byte(destKey), sk, nil, aggregate)
}

func (l LedisEmbeddedDB) Zclear(key string) (int64, error) {
	return l.db.ZClear([]byte(key))
}

func (l LedisEmbeddedDB) Sadd(key string, members ...string) (int64, error) {
	ms := make([][]byte, len(members))
	for i, m := range members {
		ms[i] = []byte(m)
	}
	return l.db.SAdd([]byte(key), ms...)
}

func (l LedisEmbeddedDB) Smembers(key string) ([]string, error) {
	resp, err := l.db.SMembers([]byte(key))
	if err != nil {
		return nil, err
	}
	ms := make([]string, len(resp))
	for i, m := range resp {
		ms[i] = string(m)
	}

	return ms, nil
}

func (l LedisEmbeddedDB) Sclear(key string) (int64, error) {
	return l.db.SClear([]byte(key))
}

var _ DB = (*LedisEmbeddedDB)(nil)
