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

			So(len(ga.Lines), ShouldEqual, 5)
			So(ga.Lines[0].Pattern, ShouldEqual, "some/pattern/path")
			So(ga.Lines[1].Pattern, ShouldEqual, "")
			So(ga.Lines[2].Pattern, ShouldEqual, "pattern2")
			So(ga.Lines[3].Pattern, ShouldEqual, "pattern3")
			So(ga.Lines[4].Pattern, ShouldEqual, "lonely/pattern")
		})

		Convey("When putting grypt for an existing entry", func() {
			ga := ParseGitAttribs(raw)
			ga.PutGryptEntry("pattern2")

			Convey("The number of lines should not change", func() {
				So(len(ga.Lines), ShouldEqual, 5)
			})
			Convey("The existing entry should now speak of grypt", func() {
				So(ga.Lines[2].Pattern, ShouldEqual, "pattern2")
				So(string(ga.Lines[2].Raw), ShouldEqual, "pattern2 filter=grypt diff=grypt")
			})
		})

		Convey("When putting grypt for an new entry", func() {
			ga := ParseGitAttribs(raw)
			ga.PutGryptEntry("you/aint/never/seen")

			Convey("The number of lines should increment", func() {
				So(len(ga.Lines), ShouldEqual, 6)
			})
			Convey("The new entry should now speak of grypt", func() {
				So(ga.Lines[5].Pattern, ShouldEqual, "you/aint/never/seen")
				So(string(ga.Lines[5].Raw), ShouldEqual, "you/aint/never/seen filter=grypt diff=grypt")
			})
		})
	})
}
