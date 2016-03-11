package gomate

type ScorePair struct {
	Score  int64
	Member string
}

const (
	AggregateSum byte = 0
	AggregateMin byte = 1
	AggregateMax byte = 2
)

type DB interface {
	Zadd(key string, pairs ...ScorePair) error
	Zrange(key string, start int, stop int) ([]ScorePair, error)
	Zinterstore(destKey string, srcKeys []string, aggregate byte) (int64, error)
	Zclear(key string) (int64, error)
	Zkeyexists(key string) (bool, error)
	Zexpire(key string, duration int64) (int64, error)
	Sadd(key string, member ...string) (int64, error)
	Smembers(key string) ([]string, error)
	Sclear(key string) (int64, error)
}
