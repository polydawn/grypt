package main

import (
	"fmt"
	"os"
	"polydawn.net/grypt/cli"
)

var EXIT_BADARGS = 1
var EXIT_PANIC = 3

func main() {
	defer panicHandler()

	cli.Run("grypt", os.Args...)
}

func panicHandler() {
	// print only the error message (don't dump stacks).
	// unless any debug mode is on; then don't recover, because we want to dump stacks.
	if len(os.Getenv("DEBUG")) == 0 {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(EXIT_PANIC)
		}
	}
}
