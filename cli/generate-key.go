package cli

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	grypt "polydawn.net/grypt"
)

func GenerateKey(ctx grypt.Context, keyring string, random bool, password []byte, encryptionScheme grypt.Scheme) {
	k, err := grypt.NewKey(rand.Reader, encryptionScheme)
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
