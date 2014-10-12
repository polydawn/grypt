package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"polydawn.net/grypt/gitutil"
	"polydawn.net/grypt/schema"
	"polydawn.net/grypt/vault"
)

func PlumbingClean(ctx gitutil.Context, in io.Reader, out io.Writer) {
	keyDir := filepath.Join(ctx.RepoDataDir, "grypt")
	k, err := schema.ReadKey(filepath.Join(keyDir, ctx.Keyring+".key"))
	if err != nil {
		panic(fmt.Errorf("error reading key: %v", err)) // you see how this is the same in every function (except for the one that was previous divergent through sheer human oversight)?  this is why that "handle errors where they occur" mantra is complete horseshit.  it leads to stupid duplication of error handling absolutely fucking everywhere, and that multiplied by time and contact with the real world leads to inconsistent error handling.  painful.
	}
	// TODO: pause for a moment... need a function to stand in here that parses the headers and then does dispatch
	if err := vault.Encrypt(ctx, in, out, k); err != nil {
		panic(fmt.Errorf("crypto error: %v", err))
	}
}

func PlumbingSmudge(ctx gitutil.Context, in io.Reader, out io.Writer) {
	keyDir := filepath.Join(ctx.RepoDataDir, "grypt")
	k, err := schema.ReadKey(filepath.Join(keyDir, ctx.Keyring+".key"))
	if err != nil {
		panic(fmt.Errorf("error reading key: %v", err))
	}
	// TODO: pause for a moment... need a function to stand in here that parses the headers and then does dispatch
	if err := grypt.Decrypt(in, out, k); err != nil {
		panic(fmt.Errorf("crypto error: %v", err))
	}
}

func PlumbingTextconv(ctx gitutil.Context, f string, out io.Writer) {
	keyDir := filepath.Join(ctx.RepoDataDir, "grypt")
	k, err := schema.ReadKey(filepath.Join(keyDir, ctx.Keyring+".key"))
	if err != nil {
		panic(fmt.Errorf("error reading key: %v", err))
	}
	file, err := os.Open(f)
	if err != nil {
		panic(fmt.Errorf("error reading textconv input file: %v", err))
	}
	defer file.Close()
	// TODO: pause for a moment... need a function to stand in here that parses the headers and then does dispatch
	if err := grypt.Decrypt(file, out, k); err != nil {
		panic(fmt.Errorf("crypto error: %v", err))
	}
}
