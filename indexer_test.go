package gomate

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGeneratePrefix(t *testing.T) {

	Convey("Given a short string", t, func() {
		doc := "ab"

		Convey("Then it returns an empty array", func() {
			p := generatePrefixes(doc)
			So(len(p), ShouldEqual, 0)
		})

	})
	Convey("Given a long string, with 6 chars", t, func() {
		doc := "abcdef"

		Convey("Then it returns an array with all 4 prefixes", func() {
			p := generatePrefixes(doc)
			So(len(p), ShouldEqual, 4)
			So(p, ShouldContain, "ab")
			So(p, ShouldContain, "abc")
			So(p, ShouldContain, "abcd")
			So(p, ShouldContain, "abcde")
		})

	})

}

type mockDB struct {
	DB
	keys map[string][]ScorePair
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

func TestIndex(t *testing.T) {
	db := &mockDB{}
	idx := NewIndexer(db)
	Convey("Given a document with a single word", t, func() {
		doc := "single"

		Convey("When I index this document", func() {
			err := idx.Index("1", doc)

			Convey("Then it should be successful", func() {
				So(err, ShouldBeNil)
			})
			Convey("And it whould add one key for each prefix", func() {
				So(db.keys, ShouldHaveLength, 5)
			})
			Convey("And the key that matches the whole word should have score 1", func() {
				So(db.keys["gomate-index:terms:single"][0].Score, ShouldEqual, 1)
				So(db.keys["gomate-index:terms:single"][0].Member, ShouldEqual, "1")
			})
		})
	})
}
