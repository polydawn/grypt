// +build windows

package main

import (
	"fmt"
	"os"
)

func readPhrase() []byte {
	fmt.Fprintln(os.Stderr, "Passphrase is not supported on this platorm")
	os.Exit(1)
	return []byte{}
}
