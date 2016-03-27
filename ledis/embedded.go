package ledis

import (
	"github.com/deluan/gomate"
	"github.com/siddontang/ledisdb/ledis"
)

type EmbeddedDB struct {
	db *ledis.DB
}

func NewEmbeddedDB(db *ledis.DB) gomate.DB {
	return &EmbeddedDB{db: db}
}

func (l *EmbeddedDB) Zadd(key string, pairs ...gomate.ScorePair) error {
	ps := make([]ledis.ScorePair, len(pairs))
	for i, p := range pairs {
		ps[i] = ledis.ScorePair{Score: p.Score, Member: []byte(p.Member)}
	}
	_, err := l.db.ZAdd([]byte(key), ps...)
	return err
}

func (l *EmbeddedDB) Zrem(key string, members ...string) error {
	ms := make([][]byte, len(members))
	for i, k := range members {
		ms[i] = []byte(k)
	}
	_, err := l.db.ZRem([]byte(key), ms...)
	return err
}

func (l *EmbeddedDB) Zrange(key string, start int, stop int) ([]gomate.ScorePair, error) {
	res, err := l.db.ZRange([]byte(key), start, stop)
	if err != nil {
		return nil, err
	}
	ps := make([]gomate.ScorePair, len(res))
	for i, p := range res {
		ps[i] = gomate.ScorePair{Score: p.Score, Member: string(p.Member)}
	}
	return ps, nil
}

func (l *EmbeddedDB) Zinterstore(destKey string, srcKeys []string, aggregate byte) (int64, error) {
	sk := make([][]byte, len(srcKeys))
	for i, k := range srcKeys {
		sk[i] = []byte(k)
	}
	return l.db.ZInterStore([]byte(destKey), sk, nil, aggregate)
}

func (l *EmbeddedDB) Zclear(key string) (int64, error) {
	return l.db.ZClear([]byte(key))
}

func (l *EmbeddedDB) Zkeyexists(key string) (bool, error) {
	resp, err := l.db.ZKeyExists([]byte(key))
	return resp == 1, err
}

func (l *EmbeddedDB) Zexpire(key string, duration int64) (int64, error) {
	return l.db.ZExpire([]byte(key), duration)
}

func (l *EmbeddedDB) Sadd(key string, members ...string) (int64, error) {
	ms := make([][]byte, len(members))
	for i, m := range members {
		ms[i] = []byte(m)
	}
	return l.db.SAdd([]byte(key), ms...)
}

func (l *EmbeddedDB) Smembers(key string) ([]string, error) {
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

func (l *EmbeddedDB) Sclear(key string) (int64, error) {
	return l.db.SClear([]byte(key))
}

func (l *EmbeddedDB) Hset(key, field, value string) (int64, error) {
	return l.db.HSet([]byte(key), []byte(field), []byte(value))
}

func (l *EmbeddedDB) Hclear(key string) (int64, error) {
	return l.db.HClear([]byte(key))
}

func (l *EmbeddedDB) Hmget(key string, fields ...string) ([]string, error) {
	fs := make([][]byte, len(fields))
	for i, f := range fields {
		fs[i] = []byte(f)
	}
	resp, err := l.db.HMget([]byte(key), fs...)
	if err != nil {
		return nil, err
	}
	values := make([]string, len(resp))
	for i, m := range resp {
		values[i] = string(m)
	}

	return values, nil
}

var _ gomate.DB = (*EmbeddedDB)(nil)
