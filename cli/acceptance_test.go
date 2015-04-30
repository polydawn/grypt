package cli

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"polydawn.net/grypt/gitutil"
	"polydawn.net/grypt/testutil"
)

//var git = gosh.Sh("git")

func TestGenerateKey(t *testing.T) {
	testutil.Hideme(func() {
		Convey("Given a new git repo with at least one commit", t, func() {
			git("init")()
			git("commit", "--allow-empty", "-m", "initial commit")()

			Convey("There should be no evidence of grypt yet", func() {
				_, err := os.Stat(".git/grypt/")
				exists := err == nil || !os.IsNotExist(err)
				So(exists, ShouldBeFalse)
			})

			Convey("When 'grypt generate-key' is called", func() {
				Run(
					"irrelephant",
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
	// gpath := strings.Join([]string{gdir, os.Getenv("PATH")}, string(os.PathListSeparator))
	// turns out trying to set $PATH doesn't work because git makes itself a fresh shell before calling filter commands.
	//git := git(gosh.Env{"PATH": gpath})
	//os.Setenv("PATH", gpath)

	Convey("Given a new repo", t,
		testutil.WithTmpdir(func() {
			git("init")()
			git("commit", "--allow-empty", "-m", "initial commit")()

			Convey("When 'grypt keep-secret' is called with no args", func() {
				Convey("We should panic", func() {
					So(func() {
						Run(
							"irrelephant",
							"grypt",
							"keep-secret",
						)
					}, ShouldPanicWith, "which files should we keep secret?")
				})
			})

			Convey("When 'grypt keep-secret not-a-file' is called", func() {
				Convey("We should panic", func() {
					So(func() {
						Run(
							"irrelephant",
							"grypt",
							"keep-secret",
							"not-a-file",
						)
					}, ShouldPanicWith, "stat not-a-file: no such file or directory")
					// this is not the greatest error message...
				})
			})

			Convey("When 'grypt keep-secret some-dir' is called", func() {
				So(os.Mkdir("some-dir", 0755), ShouldBeNil)
				SkipConvey("We should panic", func() {
					So(func() {
						Run(
							"irrelephant",
							"grypt",
							"keep-secret",
							"some-dir",
						)
					}, ShouldPanic)
					// borken
				})
			})
		}),
	)

	Convey("Given a new repo with grypt already past generate-key", t,
		testutil.WithTmpdir(func() {

			git("init")()
			git("commit", "--allow-empty", "-m", "initial commit")()
			Run(
				"irrelephant",
				"grypt",
				"generate-key",
				"--password", "asdf", // do not want interactive prompt to be hit in tests
			)

			Convey("When 'grypt keep-secret shadowfile' is called", func() {
				err := ioutil.WriteFile("shadowfile", []byte("cleartext"), 0644)
				So(err, ShouldBeNil)
				Run(
					filepath.Join(gdir, "grypt"),
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
					So(git("status", "--porcelain").Output(), ShouldEqual, "A  .gitattributes\nA  shadowfile\n")
				})

				Convey("The raw staged file should show the ciphertext", func() {
					// staged diff should have our serial ciphertext -- the inspection we're using here is low level enough that it does not give the diff filter a chance to run
					stagedLines := strings.Split(string(gitutil.ListStagedFileContents()["shadowfile"]), "\n")

					So(stagedLines[0], ShouldEqual, "-----BEGIN GRYPT CIPHERTEXT HEADER-----")
					So(git("diff", "--raw").Output(), ShouldEqual, "")
				})

				Convey("If I nuke the gitattributes on disk", func() {
					os.Remove(".gitattributes")

					// '--no-ext-diff' or '--no-textconv' might be alternative ways to test this

					Convey("The secret file should still be staged", func() {
						So(git("status", "--porcelain").Output(), ShouldEqual, "AD .gitattributes\nA  shadowfile\n")
					})

					Convey("The dircache should contain the ciphertext", func() {
						stagedLines := strings.Split(string(gitutil.ListStagedFileContents()["shadowfile"]), "\n")
						So(stagedLines[0], ShouldEqual, "-----BEGIN GRYPT CIPHERTEXT HEADER-----")
					})

					Convey("The diff should be empty", func() {
						So(git("diff", "--raw", "shadowfile").Output(), ShouldEqual, "")
					})

					Convey("The working tree should show the cleartext", func() {

					})
				})

				Convey("If I nuke the gitattributes on disk and dircache", func() {
					os.Remove(".gitattributes")
					git("rm", ".gitattributes")()

					// NOW things are different.  git will use the staged gitattributes if there isn't one on disk!

					// all of these assertions describe what happens when we have a secret staged but git isn't configured to handle it; might be able to use same assertions for both removed gitattributes or .git/config.
					// and actually, the form of these adapted to post-commit... should look quite a lot like when there's no key or the lock subcommand is called

					Convey("The secret file is still staged and now considered modified", func() {
						So(git("status", "--porcelain").Output(), ShouldEqual, "AM shadowfile\n")
					})

					Convey("The dircache should contain the ciphertext", func() {
						stagedLines := strings.Split(string(gitutil.ListStagedFileContents()["shadowfile"]), "\n")
						So(stagedLines[0], ShouldEqual, "-----BEGIN GRYPT CIPHERTEXT HEADER-----")
					})

					Convey("The diff should be dirty", func() {
						// Since we still have the cleartext on disk, and git no longer knows how to filter things, git sees this as a "change".
						// FUTURE: Contemplate: would it make sense to install a pre-commit hook that looks for files covered by a grypt filter attribute, and abort the commit if there's a staged change to that file that doesn't look like ciphertext?
						So(git("diff", "--raw", "shadowfile").Output(), ShouldNotEqual, "")
					})

					Convey("The working tree should show the cleartext", func() {
						worktreeBytes, _ := ioutil.ReadFile("shadowfile")
						worktreeLines := strings.Split(string(worktreeBytes), "\n")
						So(worktreeLines[0], ShouldEqual, "cleartext")
					})

					Convey("After I perform a checkout", func() {
						git("checkout", ".")()

						Convey("The diff should be empty", func() {
							So(git("diff", "--raw", "shadowfile").Output(), ShouldEqual, "")
						})
						Convey("The working tree should show the ciphertext", func() {
							worktreeBytes, _ := ioutil.ReadFile("shadowfile")
							worktreeLines := strings.Split(string(worktreeBytes), "\n")
							So(worktreeLines[0], ShouldEqual, "-----BEGIN GRYPT CIPHERTEXT HEADER-----")
						})
					})
				})
			})

			Convey("When 'grypt keep-secret already-staged-file' is called", func() {
				Convey("We should ...?", nil)
				// error?
				// replace the staged content with the encrypted version?
				// depends on whether previously committed?  warn if so?
				// remember that our security model is based on the presumption you still have faith in your own local disk; you just don't want to push anything too interesting.
			})

			Convey("When 'grypt keep-secret already-committed-file' is called", func() {
				Convey("We should ...?", nil)
				// error?  warn (loudly)?
				// replace the staged content with the encrypted version?
				// how do diffs look when we add a gitattribute for a path that previously existed but didn't have the attribute?
			})
		}),
	)

	Convey("Given a clone of a repo with secrets", t,
		testutil.WithTmpdir(func() {
			So(os.Mkdir("upstream", 0755), ShouldBeNil)
			So(os.Chdir("upstream"), ShouldBeNil)

			git("init")()
			git("commit", "--allow-empty", "-m", "initial commit")()
			Run(
				"irrelephant",
				"grypt",
				"generate-key",
				"--password", "asdf", // do not want interactive prompt to be hit in tests
			)
			So(ioutil.WriteFile("shadowfile", []byte("cleartext"), 0644), ShouldBeNil)
			Run(
				filepath.Join(gdir, "grypt"),
				"grypt",
				"keep-secret",
				"shadowfile",
			)

			git("commit", "-m", "keeping a secret")()

			So(os.Chdir(".."), ShouldBeNil)

			git("clone", "./upstream", "consumer")() // this not panicking is important, for starters.
			So(os.Chdir("consumer"), ShouldBeNil)

			Convey("", func() {
				logLines := strings.Split(git("log", "--oneline").Output(), "\n")
				So(len(logLines), ShouldEqual, 3) // the empty initial commit, and the commit where we added the secret (and a trailing line)

				Convey("No files are staged or modified", func() {
					So(git("status", "--porcelain").Output(), ShouldEqual, "")
				})

				Convey("The dircache should contain nothing", func() {
					stagedLines := strings.Split(string(gitutil.ListStagedFileContents()["shadowfile"]), "\n")
					So(len(stagedLines), ShouldEqual, 1)
					So(stagedLines[0], ShouldEqual, "")
				})

				Convey("The diff should be clean", func() {
					So(git("diff", "--raw", "shadowfile").Output(), ShouldEqual, "")
				})

				Convey("The working tree should show the ciphertext", func() {
					worktreeBytes, _ := ioutil.ReadFile("shadowfile")
					worktreeLines := strings.Split(string(worktreeBytes), "\n")
					So(worktreeLines[0], ShouldEqual, "-----BEGIN GRYPT CIPHERTEXT HEADER-----")
				})

				Convey("After I perform a checkout", func() {
					// i mean, this certainly shouldn't change anything.  but let's just check.
					git("checkout", ".")()

					Convey("The diff should be empty", func() {
						So(git("diff", "--raw", "shadowfile").Output(), ShouldEqual, "")
					})
					Convey("The working tree should show the ciphertext", func() {
						worktreeBytes, _ := ioutil.ReadFile("shadowfile")
						worktreeLines := strings.Split(string(worktreeBytes), "\n")
						So(worktreeLines[0], ShouldEqual, "-----BEGIN GRYPT CIPHERTEXT HEADER-----")
					})
				})
			})
		}),
	)
}
