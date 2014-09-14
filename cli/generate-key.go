package cli

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"code.google.com/p/go.crypto/hkdf"
	grypt "polydawn.net/grypt"
)

func GenerateKey(ctx grypt.Context, keyring string, random bool, password []byte, encryptionScheme grypt.Scheme) {
	var k grypt.Key
	var err error
	if random {
		k, err = grypt.NewKey(rand.Reader, encryptionScheme)
	} else {
		hkdf := hkdf.New(sha256.New, password, nil, nil)
		k, err = grypt.NewKey(hkdf, encryptionScheme)
	}
	if err != nil {
		panic(fmt.Errorf("failure generating key: %v", err))
	}

	keyDir := filepath.Join(ctx.RepoDataDir, "grypt")
	os.Mkdir(keyDir, 0700)
	keyPath := filepath.Join(keyDir, keyring+".key") // TODO: should whitelist patterns for 'keyring'

	err = grypt.WriteKey(keyPath, k)
	if err != nil {
		panic(fmt.Errorf("failure saving key: %v", err))
	}
}
