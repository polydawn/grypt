package cli

import (
	"os"
	"polydawn.net/pogo/gosh"
	"polydawn.net/grypt/testutil"
	"testing"
)

func TestWhereami(t *testing.T) {
	pwd, _ := os.Getwd()
	println(pwd)
	gosh.Sh("pwd")(gosh.DefaultIO)()
}

func TestWhereami2(t *testing.T) {
	testutil.Hideme(func() {
		pwd, _ := os.Getwd()
		println(pwd)
		gosh.Sh("pwd")(gosh.DefaultIO)()
	})
}
