package gomate

import (
	"sync"
	"testing"

	"github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	_ledisInstance *ledis.Ledis
	_dbInstance    *ledis.DB
	once           sync.Once
)

func db() *ledis.DB {
	once.Do(func() {
		config := config.NewConfigDefault()
		config.DataDir = "tmp-gomate-test"
		l, _ := ledis.Open(config)
		instance, err := l.Select(0)
		if err != nil {
			panic(err)
		}
		_ledisInstance = l
		_dbInstance = instance
	})
	return _dbInstance
}

func dropDb() {
	db()
	_ledisInstance.FlushAll()
}

func TestSearcher(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	db := NewLedisEmbeddedDB(db())
	idx := NewIndexer(db)
	s := NewSearcher(db)
	SkipConvey("Given an index with a single word document", t, func() {
		doc := "single"
		idx.Index("1", doc)
		Convey("When I search for that word", func() {
			res, err := s.Search("single")
			Convey("Then I get the id for that doc", func() {
				So(err, ShouldBeNil)
				So(res[0], ShouldEqual, "1")
			})
		})
		Reset(func() {
			dropDb()
		})
	})

}
