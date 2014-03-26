// +build !windows

package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"code.google.com/p/go.crypto/ssh/terminal"
)

func readPhrase() []byte {
	h := sha256.New()
	fmt.Fprintf(os.Stderr, "passphrase: ")
	p, err := terminal.ReadPassword(0)
	if err != nil {
		panic(err)
	}
	io.Copy(h, bytes.NewBuffer(p))
	fmt.Println()
	return h.Sum(nil)
}
