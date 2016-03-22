package gomate

import "strings"

type mockDB struct {
	DB
	keys map[string]map[string]ScorePair
	set  map[string]bool
	hash map[string]string
}

func NewMockDB() *mockDB {
	db := &mockDB{}
	db.keys = make(map[string]map[string]ScorePair)
	db.set = make(map[string]bool)
	db.hash = make(map[string]string)
	return db
}

func (db *mockDB) Zadd(key string, pairs ...ScorePair) error {
	if db.keys[key] == nil {
		db.keys[key] = make(map[string]ScorePair)
	}
	for _, p := range pairs {
		db.keys[key][p.Member] = p
	}
	return nil
}

func (db *mockDB) Zrem(key string, members ...string) error {
	set := db.keys[key]
	for _, m := range members {
		delete(set, m)
	}
	return nil
}

func (db *mockDB) Sadd(key string, member ...string) (int64, error) {
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

func (db *mockDB) Hset(key, field, value string) (int64, error) {
	_, exists := db.hash[field]
	db.hash[field] = value

	if exists {
		return 0, nil
	}
	return 1, nil
}

func (db *mockDB) Hmget(key string, fields ...string) ([]string, error) {
	var resp []string
	for _, f := range fields {
		resp = append(resp, db.hash[f])
	}
	return resp, nil
}

func (db *mockDB) Hclear(key string) (int64, error) {
	num := len(db.hash)
	db.hash = make(map[string]string)
	return int64(num), nil
}

func (db *mockDB) IsEmpty() bool {
	return len(db.keys) == 0 && len(db.set) == 0
}
