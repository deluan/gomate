package gomate

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGeneratePrefix(t *testing.T) {

	Convey("Given a one letter string", t, func() {
		doc := "a"

		Convey("Then it returns an empty array", func() {
			p := generatePrefixes(doc)
			So(len(p), ShouldEqual, 0)
		})

	})
	Convey("Given a two letters string", t, func() {
		doc := "ab"

		Convey("Then it returns an array with only the first letter", func() {
			p := generatePrefixes(doc)
			So(len(p), ShouldEqual, 1)
			So(p, ShouldContain, "a")
		})

	})
	Convey("Given a long string, with 6 chars", t, func() {
		doc := "abcdef"

		Convey("Then it returns an array with all 5 prefixes", func() {
			p := generatePrefixes(doc)
			So(len(p), ShouldEqual, 5)
			So(p, ShouldContain, "a")
			So(p, ShouldContain, "ab")
			So(p, ShouldContain, "abc")
			So(p, ShouldContain, "abcd")
			So(p, ShouldContain, "abcde")
		})

	})

}

func TestIndex(t *testing.T) {
	db := NewMockDB()
	idx := NewIndexer(db)
	Convey("Given a document with a single word", t, func() {
		doc := "single"

		Convey("When I index this document", func() {
			err := idx.Index("1", doc)

			Convey("Then it should be successful", func() {
				So(err, ShouldBeNil)
			})
			Convey("And it whould add one key for each prefix", func() {
				total := 5 + 1 // 1 extra for the keychain
				So(db.keys, ShouldHaveLength, total)
			})
			Convey("And the key that matches the whole word should be the first result", func() {
				So(db.keys["gomate-index:terms:single"]["1"].Score, ShouldEqual, 0)
			})
			Convey("And it should collect one key for each prefix", func() {
				So(db.set, ShouldHaveLength, 5+2)
			})
			Convey("And when I call Clear, it deletes all keys from the keychain", func() {
				So(db.IsEmpty(), ShouldBeFalse)
				idx.Clear()
				So(db.IsEmpty(), ShouldBeTrue)
			})
		})
	})
}
