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
		ga := []byte(strings.Join([]string{
			"some/pattern/path k1=val1 k2=val2",
		}, "\n"))

		Convey("There should be a new entry", func() {
			// TODO finish
			So(ga, ShouldBeTrue)
		})
	})
}
