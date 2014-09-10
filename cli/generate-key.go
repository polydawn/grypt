package cli

import (
	"crypto/rand"
	"fmt"
	grypt "polydawn.net/grypt"
)

func GenerateKey(keyring string, random bool, password []byte, encryptionScheme grypt.Scheme) {
	k, err := grypt.NewKey(rand.Reader, encryptionScheme)
	if err != nil {
		panic(fmt.Errorf("failure generating key: %v", err))
	}
	err = grypt.WriteKey("hardcoded_fixme", k)
	if err != nil {
		panic(fmt.Errorf("failure saving key: %v", err))
	}
}
