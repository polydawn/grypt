package cli

import (
	"fmt"
	"os"
	"path/filepath"
	grypt "polydawn.net/grypt"
)

// TODO: all these currently consume hardcoded Stdin and Stdout.  this should probably be parameterized

func PlumbingClean(ctx grypt.Context, keyring string) {
	keyDir := filepath.Join(ctx.RepoDataDir, "grypt")
	k, err := grypt.ReadKey(filepath.Join(keyDir, keyring+".key"))
	if err != nil {
		panic(fmt.Errorf("error reading key: %v", err)) // you see how this is the same in every function (except for the one that was previous divergent through sheer human oversight)?  this is why that "handle errors where they occur" mantra is complete horseshit.  it leads to stupid duplication of error handling absolutely fucking everywhere, and that multiplied by time and contact with the real world leads to inconsistent error handling.  painful.
	}
	if err := grypt.Encrypt(os.Stdin, os.Stdout, k); err != nil {
		panic(fmt.Errorf("crypto error: %v", err))
	}
}

func PlumbingSmudge(ctx grypt.Context, keyring string) {
	keyDir := filepath.Join(ctx.RepoDataDir, "grypt")
	k, err := grypt.ReadKey(filepath.Join(keyDir, keyring+".key"))
	if err != nil {
		panic(fmt.Errorf("error reading key: %v", err))
	}
	if err := grypt.Decrypt(os.Stdin, os.Stdout, k); err != nil {
		panic(fmt.Errorf("crypto error: %v", err))
	}
}

func PlumbingTextconv(ctx grypt.Context, keyring string, f string) {
	keyDir := filepath.Join(ctx.RepoDataDir, "grypt")
	k, err := grypt.ReadKey(filepath.Join(keyDir, keyring+".key"))
	if err != nil {
		panic(fmt.Errorf("error reading key: %v", err))
	}
	file, err := os.Open(f)
	if err != nil {
		panic(fmt.Errorf("error reading textconv input file: %v", err))
	}
	defer file.Close()
	if err := grypt.Decrypt(file, os.Stdout, k); err != nil {
		panic(fmt.Errorf("crypto error: %v", err))
	}
}
