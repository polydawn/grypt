package testutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"polydawn.net/pogo/gosh"
)

/*
	Produces a grypt binary by exec'ing `go build`.
	Returns a string representing a directory you can put on your $PATH in order to use `grypt`.

	No other binaries will be in the returned directory (so you should be able to put it on the *front* of your $PATH,
	as indeed you'll need to in order to use this sensibly if your environment already has a real grypt command available).
*/
func BuildGrypt() string {
	// TODO: i don't really know how I intend to implement cleanup of this thing.
	// maybe we'll end up with another function chainer like WithTmpdir() that lets us just bind it to goconvey context.  but there's also no earthly reason we'll want to tolerate running the compiler again for every test.

	// find the project root dir.  go tests run with a cwd that is the directory of the package.  so we should be able to look up until we find an indicator that we're in the project root.
	// the indicator chosen is the goad script, because fuck it, it makes as little sense as anything else.  if there's an alternative where i can ask the go testing package where i *wish* my cwd was, that'd be great.
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	projectBase := cwd
	for {
		projectBase = filepath.Dir(projectBase)
		_, err := os.Stat(filepath.Join(projectBase, "goad"))
		exists := err == nil || !os.IsNotExist(err)
		if exists {
			break // gotcha
		}
		if projectBase == string(filepath.Separator) {
			panic(fmt.Errorf("reached root without finding project base dir when starting from %s", cwd))
		}
	}

	// make a tmpdir where we can output the binary
	tmpBase := os.Getenv("TMPDIR")
	if len(tmpBase) == 0 {
		tmpBase = os.TempDir()
	}
	tmpBase = filepath.Join(tmpBase, "grypt-build")
	err = os.MkdirAll(tmpBase, 0755)
	if err != nil {
		panic(err)
	}
	tmpdir, err := ioutil.TempDir(tmpBase, "")
	if err != nil {
		panic(err)
	}

	// exec `go build` with a GOPATH based on the project dir and have it put a grypt exectuable in our tempdir.
	// we can leave stdin/stdout/stderr connected to the compiler exec because it's silent, unless something goes wrong in which case we do quite want to see it.
	g := gosh.Sh("go")(gosh.Env{"GOPATH": filepath.Join(projectBase, ".gopath")})
	g("build", "-o", filepath.Join(tmpdir, "grypt"), "polydawn.net/grypt/main")(gosh.DefaultIO)()

	return tmpdir
}
