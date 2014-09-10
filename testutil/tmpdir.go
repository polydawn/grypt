package testutil

import (
	"io/ioutil"
	"path/filepath"
	"os"
)

func Hideme(fn func()) {
	retreat, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	defer os.Chdir(retreat)

	tmpBase := os.Getenv("TMPDIR")
	if len(tmpBase) == 0 {
		tmpBase = os.TempDir()
	}
	err = os.Chdir(tmpBase)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("grypt-test", 0755)
	if err != nil {
		panic(err)
	}
	tmpdir, err := ioutil.TempDir("grypt-test", "")
	if err != nil {
		panic(err)
	}
	tmpdir, err = filepath.Abs(tmpdir)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpdir)
	err = os.Chdir(tmpdir)
	if err != nil {
		panic(err)
	}

	fn()
}
