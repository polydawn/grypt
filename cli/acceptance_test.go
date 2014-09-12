package cli

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"polydawn.net/grypt/cli"
	"polydawn.net/grypt/gitutil"
	"polydawn.net/grypt/testutil"
	"polydawn.net/pogo/gosh"
	"testing"
)

var git = gosh.Sh("git")

func TestGenerateKey(t *testing.T) {
	testutil.Hideme(func() {
		Convey("Given a new git repo with at least one commit", t, func() {
			git("init")()
			git("commit")("--allow-empty", "-m", "initial commit")()

			Convey("There should be no evidence of grypt yet", func() {
				_, err := os.Stat(".git/grypt/")
				exists := err == nil || !os.IsNotExist(err)
				So(exists, ShouldBeFalse)
			})

			Convey("When 'grypt generate-key' is called", func() {
				cli.Run(
					"grypt",
					"generate-key",
					"--password", "asdf", // do not want interactive prompt to be hit in tests
				)

				Convey("We should get a key file in the git data dir", func() {

					// ridiculously verbose golang way to check if a file exists
					// we should wrap up this up in a custom assertion, because we're gonna be doing it a lot: https://github.com/smartystreets/goconvey/wiki/Custom-Assertions
					_, err := os.Stat(".git/grypt/default.key")
					exists := err == nil || !os.IsNotExist(err)

					So(exists, ShouldBeTrue)
				})
			})
		})
	})
}

func TestKeepSecret(t *testing.T) {
	testutil.Hideme(func() {
		Convey("Given a new repo with grypt already past generate-key", t, func() {
			git("init")()
			git("commit")("--allow-empty", "-m", "initial commit")()
			cli.Run(
				"grypt",
				"generate-key",
				"--password", "asdf", // do not want interactive prompt to be hit in tests
			)

			Convey("When 'grypt keep-secret shadowfile' is called", func() {
				cli.Run(
					"grypt",
					"keep-secret",
					"shadowfile",
				)

				Convey("We should see gitattributes", func() {
					_, err := os.Stat(".gitattributes")
					exists := err == nil || !os.IsNotExist(err)
					So(exists, ShouldBeTrue)

					ga := gitutil.ReadGitAttribsFile(".gitattributes")
					So(len(ga.Lines), ShouldEqual, 1)
					So(ga.Lines[0].Pattern, ShouldEqual, "shadowfile")
				})
			})
		})
	})
}
