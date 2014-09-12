package testutil

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"path/filepath"
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

/*
	Decorates a goconvey test with a tmpdir.

	See also https://github.com/smartystreets/goconvey/wiki/Decorating-tests-to-provide-common-logic
*/
func WithTmpdir(fn func()) func() {
	// yes, this function is entirely ths same as 'Hideme()' but with another layer of closure and a goconvey 'Reset()' instead of a defer.

	return func() {
		retreat, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		Reset(func() {
			os.Chdir(retreat)
		})

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
		Reset(func() {
			os.RemoveAll(tmpdir)
		})
		err = os.Chdir(tmpdir)
		if err != nil {
			panic(err)
		}

		fn()
	}
}
