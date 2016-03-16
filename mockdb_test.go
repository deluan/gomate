package gomate

import "strings"

type mockDB struct {
	DB
	keys map[string][]ScorePair
	set  map[string]bool
}

func (db *mockDB) Zadd(key string, pairs ...ScorePair) error {
	if db.keys == nil {
		db.keys = make(map[string][]ScorePair)
	}
	if db.keys[key] == nil {
		db.keys[key] = make([]ScorePair, 0)
	}
	db.keys[key] = append(db.keys[key], pairs...)
	return nil
}

func (db *mockDB) Sadd(key string, member ...string) (int64, error) {
	if db.set == nil {
		db.set = make(map[string]bool)
	}
	for _, m := range member {
		db.set[m] = true
	}
	return int64(len(member)), nil
}

func (db *mockDB) Smembers(key string) ([]string, error) {
	m := make([]string, 0, len(db.set))
	for k := range db.set {
		m = append(m, k)
	}
	return m, nil
}

func (db *mockDB) Sclear(key string) (int64, error) {
	if strings.HasSuffix(key, KeyChainSuffix) {
		db.set = make(map[string]bool)
	} else {
		delete(db.set, key)
	}
	return 1, nil
}

func (db *mockDB) Zclear(key string) (int64, error) {
	delete(db.keys, key)
	return 1, nil
}

func (db *mockDB) IsEmpty() bool {
	return len(db.keys) == 0 && len(db.set) == 0
}
