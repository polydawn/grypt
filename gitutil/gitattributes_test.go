package gitutil

import (
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

// var git = gosh.Sh("git")
// i was going to propose that we test if git can read back our files, but that's really remarkably hard to do
// the git-check-attr command is severely braindamaged about controlling its inputs: specifically, you can't actually give it a list of where to look for gitattributes files so that the test can actually be isolated from your system.

func TestGitattributesFiltering(t *testing.T) {
	Convey("Given some gitattributes string", t, func() {
		raw := []byte(strings.Join([]string{
			"some/pattern/path k1=val1 k2=val2",
			"",
			"pattern2\t\tka=val1",
			"pattern3 kb=val1",
			"lonely/pattern",
		}, "\n"))

		Convey("The entries should be parsible", func() {
			ga := ParseGitAttribs(raw)

			So(len(ga.lines), ShouldEqual, 5)
			So(ga.lines[0].pattern, ShouldEqual, "some/pattern/path")
			So(ga.lines[1].pattern, ShouldEqual, "")
			So(ga.lines[2].pattern, ShouldEqual, "pattern2")
			So(ga.lines[3].pattern, ShouldEqual, "pattern3")
			So(ga.lines[4].pattern, ShouldEqual, "lonely/pattern")
		})

		Convey("When putting grypt for an existing entry", func() {
			ga := ParseGitAttribs(raw)
			ga.PutGryptEntry("pattern2")

			Convey("The number of lines should not change", func() {
				So(len(ga.lines), ShouldEqual, 5)
			})
			Convey("The existing entry should now speak of grypt", func() {
				So(ga.lines[2].pattern, ShouldEqual, "pattern2")
				So(string(ga.lines[2].line), ShouldEqual, "pattern2 filter=grypt diff=grypt")
			})
		})

		Convey("When putting grypt for an new entry", func() {
			ga := ParseGitAttribs(raw)
			ga.PutGryptEntry("you/aint/never/seen")

			Convey("The number of lines should increment", func() {
				So(len(ga.lines), ShouldEqual, 6)
			})
			Convey("The existing entry should now speak of grypt", func() {
				So(ga.lines[5].pattern, ShouldEqual, "you/aint/never/seen")
				So(string(ga.lines[5].line), ShouldEqual, "you/aint/never/seen filter=grypt diff=grypt")
			})
		})
	})
}
