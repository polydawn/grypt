package testutil

import (
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

	// get the cwd, which we assume to still be the root of the project, or else you can fuck right off
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
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

	g := gosh.Sh("go")(gosh.Env{"GOPATH": filepath.Join(cwd, ".gopath")})
	g("build", "-o", filepath.Join(tmpdir, "grypt"))(gosh.DefaultIO)()

	return tmpdir
}
