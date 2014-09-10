package cli

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"polydawn.net/grypt/cli"
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

			Convey("When 'grypt generate-key' is called", func() {
				cli.Run(
					"generate-key",
					"--passowrd", "asdf", // do not want interactive prompt to be hit in tests
				)

				Convey("We should get a key file in the git data dir", func() {

					// ridiculously verbose golang way to check if a file exists
					// we should wrap up this up in a custom assertion, because we're gonna be doing it a lot: https://github.com/smartystreets/goconvey/wiki/Custom-Assertions
					_, err := os.Stat(".git/grypt/default.key")
					exists := err == nil || os.IsNotExist(err)

					So(exists, ShouldBeTrue)
				})
			})
		})
	})
}
