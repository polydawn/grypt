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
	// FIXME: oh boy, how do we get the schema here?  retain the old one?  ... for the file that might not exist before, yeah, that'll work great.  Key type used to have this for a reason.  though it was quietly being serialized as a meaningless magical int.
	vault.WeaveBasket(ctx, in, out, k)
}

func PlumbingSmudge(ctx gitutil.Context, in io.Reader, out io.Writer) {
	keyDir := filepath.Join(ctx.RepoDataDir, "grypt")
	k, err := schema.ReadKey(filepath.Join(keyDir, ctx.Keyring+".key"))
	if err != nil {
		panic(fmt.Errorf("error reading key: %v", err))
	}
	if err := grypt.OpenCiphertext(in, out, k); err != nil {
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
	if err := grypt.OpenCiphertext(file, out, k); err != nil {
		panic(fmt.Errorf("crypto error: %v", err))
	}
}
