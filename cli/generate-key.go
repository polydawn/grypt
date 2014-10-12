package cli

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"code.google.com/p/go.crypto/hkdf"

	"polydawn.net/grypt/gitutil"
	"polydawn.net/grypt/schema"
)

func GenerateKey(ctx gitutil.Context, random bool, password []byte, encryptionScheme schema.Schema) {
	var k schema.Key
	var err error
	if random {
		k, err = encryptionScheme.NewKey(rand.Reader)
	} else {
		hkdf := hkdf.New(sha256.New, password, nil, nil)
		k, err = encryptionScheme.NewKey(hkdf)
	}
	if err != nil {
		panic(fmt.Errorf("failure generating key: %v", err))
	}

	keyDir := filepath.Join(ctx.RepoDataDir, "grypt")
	os.Mkdir(keyDir, 0700)
	keyPath := filepath.Join(keyDir, ctx.Keyring+".key") // TODO: should whitelist patterns for 'keyring'

	err = schema.WriteKey(keyPath, k)
	if err != nil {
		panic(fmt.Errorf("failure saving key: %v", err))
	}
}
