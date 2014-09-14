package cli

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"polydawn.net/grypt/cli"
	"polydawn.net/grypt/gitutil"
	"polydawn.net/grypt/testutil"
	"polydawn.net/pogo/gosh"
	"strings"
	"testing"
)

//var git = gosh.Sh("git")

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
	gdir := testutil.BuildGrypt()
	gpath := strings.Join([]string{gdir, os.Getenv("PATH")}, string(os.PathListSeparator))
	git := git(gosh.Env{"PATH": gpath})

	// gosh.Sh("sh")(gosh.Env{"PATH": gpath})("-c", "echo $PATH; grypt --hallo")(gosh.DefaultIO)()

	Convey("Given a new repo with grypt already past generate-key", t,
		testutil.WithTmpdir(func() {
			os.Setenv("PATH", gpath)

			// println("------")
			// gosh.Sh("sh")("-c", "echo $PATH")(gosh.DefaultIO)()

			git("init")()
			git("commit")("--allow-empty", "-m", "initial commit")()
			cli.Run(
				"grypt",
				"generate-key",
				"--password", "asdf", // do not want interactive prompt to be hit in tests
			)

			Convey("When 'grypt keep-secret shadowfile' is called", func() {
				err := ioutil.WriteFile("shadowfile", []byte("cleartext"), 0644)
				So(err, ShouldBeNil)
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

				Convey("The secret file should be staged", func() {
					// TODO
					So(git("status", "--porcelain").Output(), ShouldEqual, "A  .gitattributes\nA  shadowfile\n")
				})

				Convey("The diff should show the cleartext", func() {
					// TODO: this might actually be insane, because it relies on exec roundtripping through another grypt process.

					// staged diff should have our bby
					// TODO: this test might actually be jumping the gun a little... i don't think the inspection we're using here is going to give the diff filter a chance to run
					// ...so it should actually probably see the ciphertext.  but we're not done connecting those parts yet, so we'll have to review this again later.
					So(string(gitutil.ListStagedFileContents()["shadowfile"]), ShouldEqual, "cleartext")
					SkipSo(1, ShouldEqual, 2) // dummy line to highlight in the report we're not done here
					// no unstaged changes should be around
					So(git("diff", "--raw").Output(), ShouldEqual, "")
				})

				Convey("If I nuke the gitattributes", func() {
					//TODO

					// '--no-ext-diff' or '--no-textconv' might be alternative ways to test this

					Convey("The diff should show the ciphertext", nil)
				})
			})
		}),
	)
}
