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
		config.DBName = "memory"
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
	Convey("Given an index with a single small word document", t, func() {
		doc := "ab"
		idx.Index("1", doc)
		Convey("When I search for that word", func() {
			res, err := s.Search("ab", 0, -1)
			Convey("Then I get the id for that doc", func() {
				So(err, ShouldBeNil)
				So(res, ShouldHaveLength, 1)
				So(res[0], ShouldEqual, "1")
			})
		})
		Reset(func() {
			dropDb()
		})
	})

	Convey("Given an index with a single word document", t, func() {
		doc := "single"
		idx.Index("2", doc)
		Convey("When I search for that specific word", func() {
			res, err := s.Search("single", 0, -1)
			Convey("Then I get the id for that doc", func() {
				So(err, ShouldBeNil)
				So(res, ShouldHaveLength, 1)
				So(res[0], ShouldEqual, "2")
			})
		})
		Convey("When I search for the first 2 chars of that word", func() {
			res, err := s.Search("si", 0, -1)
			Convey("Then I get the id for that doc", func() {
				So(err, ShouldBeNil)
				So(res, ShouldHaveLength, 1)
				So(res[0], ShouldEqual, "2")
			})
		})
		Reset(func() {
			dropDb()
		})
	})

	Convey("Given an index with a two words document", t, func() {
		doc := "joy division"
		idx.Index("3", doc)
		Convey("When I search for a non-matching", func() {
			res, err := s.Search("new", 0, -1)
			Convey("Then I get an empty result", func() {
				So(err, ShouldBeNil)
				So(res, ShouldHaveLength, 0)
			})
		})
		Convey("When I search for one of the words", func() {
			res, err := s.Search("division", 0, -1)
			Convey("Then I get the id for that doc", func() {
				So(err, ShouldBeNil)
				So(res, ShouldHaveLength, 1)
				So(res[0], ShouldEqual, "3")
			})
		})
		Convey("When I search for the first 2 chars of one of the words", func() {
			res, err := s.Search("jo", 0, -1)
			Convey("Then I get the id for that doc", func() {
				So(err, ShouldBeNil)
				So(res, ShouldHaveLength, 1)
				So(res[0], ShouldEqual, "3")
			})
		})
		Convey("When I search for the first 2 chars of each word", func() {
			res, err := s.Search("jo di", 0, -1)
			Convey("Then I get the id for that doc", func() {
				So(err, ShouldBeNil)
				So(res, ShouldHaveLength, 1)
				So(res[0], ShouldEqual, "3")
			})
		})
		Reset(func() {
			dropDb()
		})
	})
	Convey("Given an index with some documents", t, func() {
		idx.Index("1", "echo & the bunnymen")
		idx.Index("2", "erasure")
		idx.Index("3", "echoes of an era")
		Convey("When I search for one of the words", func() {
			res, err := s.Search("bunnymen", 0, -1)
			Convey("Then I get the id for that doc", func() {
				So(err, ShouldBeNil)
				So(res, ShouldHaveLength, 1)
				So(res[0], ShouldEqual, "1")
			})
		})
		Convey("When I search for a matching prefix for two docs", func() {
			res, err := s.Search("ech", 0, -1)
			Convey("Then I get the id for that both docs", func() {
				So(err, ShouldBeNil)
				So(res, ShouldHaveLength, 2)
				So(res, ShouldContain, "1")
				So(res, ShouldContain, "3")
			})
		})
		Convey("When I search for a word that matchs a whole word and a partial word for different docs", func() {
			res, err := s.Search("era", 0, -1)
			Convey("Then I get the id for that both docs", func() {
				So(err, ShouldBeNil)
				So(res, ShouldHaveLength, 2)
				So(res, ShouldContain, "2")
				So(res, ShouldContain, "3")
			})
			Convey("And I got the doc with the matching word first", func() {
				So(res[0], ShouldEqual, "3")
			})
		})
		Convey("When I search for a non-matching combination", func() {
			res, err := s.Search("bunny era", 0, -1)
			Convey("Then I get no documents back", func() {
				So(err, ShouldBeNil)
				So(res, ShouldHaveLength, 0)
			})
		})
		Convey("When I exclude one of the documents and do a search", func() {
			idx.Remove("1")
			res, err := s.Search("echo", 0, -1)
			Convey("Then the excluded doc is not returned", func() {
				So(err, ShouldBeNil)
				So(res, ShouldHaveLength, 1)
				So(res, ShouldNotContain, "1")
				So(res, ShouldContain, "3")
			})
		})
		Reset(func() {
			dropDb()
		})
	})
}
