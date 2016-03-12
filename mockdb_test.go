package gomate

type mockDB struct {
	DB
	keys map[string][]ScorePair
	kc   map[string]bool
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
	if db.kc == nil {
		db.kc = make(map[string]bool)
	}
	for _, m := range member {
		db.kc[m] = true
	}
	return int64(len(member)), nil
}

func (db *mockDB) Smembers(key string) ([]string, error) {
	m := make([]string, 0, len(db.kc))
	for k := range db.kc {
		m = append(m, k)
	}
	return m, nil
}

func (db *mockDB) Sclear(key string) (int64, error) {
	delete(db.kc, key)
	return 1, nil
}

func (db *mockDB) Zclear(key string) (int64, error) {
	delete(db.kc, key)
	return 1, nil
}
