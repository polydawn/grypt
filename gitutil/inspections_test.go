package gitutil

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"polydawn.net/grypt/testutil"
	"polydawn.net/pogo/gosh"
	"testing"
)

func TestInspectingStagedObjects(t *testing.T) {
	Convey("Given a repo with two staged new files", t,
		testutil.WithTmpdir(func() {
			git = gosh.Sh("git")
			git("init")()
			git("commit")("--allow-empty", "-m", "initial commit")()

			err := ioutil.WriteFile("alpha", []byte("111"), 0644)
			So(err, ShouldBeNil)
			err = ioutil.WriteFile("beta", []byte("222"), 0644)
			So(err, ShouldBeNil)

			git("add", "--", "alpha", "beta")()

			Convey("ListStagedFiles should return the right filenames", func() {
				stagedObjectIds := ListStagedFiles()
				So(len(stagedObjectIds), ShouldEqual, 2)
				_, ok := stagedObjectIds["alpha"]
				So(ok, ShouldBeTrue)
				_, ok = stagedObjectIds["beta"]
				So(ok, ShouldBeTrue)
			})

			Convey("ShowStagedFileContents should return the raw contents", func() {
				stagedContents := ListStagedFileContents()
				So(len(stagedContents), ShouldEqual, 2)
				So(stagedContents["alpha"], ShouldResemble, []byte("111"))
				So(stagedContents["beta"], ShouldResemble, []byte("222"))
			})
		}),
	)

	// TODO also test that we get the right (full) contents when staging a diff on top of existing file

	// TODO also test repo initialization?  though if git makes that hard, i'm very willing to give that up as a feature.  most of git itself completely sucks on a zero-commit repo.
}
