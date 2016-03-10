package gomate

type ScorePair struct {
	Score  int64
	Member string
}


type DB interface {
	Zadd(key string, pairs ...ScorePair) error
	Zrange(key string, start int, stop int) ([]ScorePair, error)
}
