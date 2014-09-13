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

	Convey("Given a repo with existing commits", t,
		testutil.WithTmpdir(func() {
			git = gosh.Sh("git")
			git("init")()
			git("commit")("--allow-empty", "-m", "initial commit")()

			err := ioutil.WriteFile("alpha", []byte("111"), 0644)
			So(err, ShouldBeNil)
			err = ioutil.WriteFile("beta", []byte("222"), 0644)
			So(err, ShouldBeNil)

			git("add", "--", "alpha", "beta")()
			git("commit")("-m", "actual content commit")()

			Convey("ListStagedFiles should be empty", func() {
				stagedObjectIds := ListStagedFiles()
				So(len(stagedObjectIds), ShouldEqual, 0)
			})

			Convey("When additional changes are made to an existing file but not staged", func() {
				err := ioutil.WriteFile("alpha", []byte("333"), 0644)
				So(err, ShouldBeNil)

				Convey("ListStagedFiles should be empty", func() {
					stagedObjectIds := ListStagedFiles()
					So(len(stagedObjectIds), ShouldEqual, 0)
				})
			})

			Convey("When additional changes are staged to an existing file", func() {
				err := ioutil.WriteFile("alpha", []byte("333"), 0644)
				So(err, ShouldBeNil)

				git("add", "--", "alpha", "beta")()

				Convey("ListStagedFiles should return the right filenames", func() {
					stagedObjectIds := ListStagedFiles()
					So(len(stagedObjectIds), ShouldEqual, 1)
					_, ok := stagedObjectIds["alpha"]
					So(ok, ShouldBeTrue)
				})

				Convey("ShowStagedFileContents should return the raw contents", func() {
					stagedContents := ListStagedFileContents()
					So(len(stagedContents), ShouldEqual, 1)
					So(stagedContents["alpha"], ShouldResemble, []byte("333"))
				})
			})
		}),
	)

	// TODO also test repo initialization?  though if git makes that hard, i'm very willing to give that up as a feature.  most of git itself completely sucks on a zero-commit repo.
}
