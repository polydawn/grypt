package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/tabwriter"

	"code.google.com/p/go.crypto/hkdf"
)

var (
	keyfile string
	exe     string

	iterations     = 4096
	attributesHelp = `Edit your .gitattributes if it's not configured already:

	secretfile filter=grypt diff=grypt
	*.secret filter=grypt diff=grypt
`
)

type (
	// Header contains information about the encryption scheme
	Header struct {
		Scheme Scheme
		MAC    []byte
	}
)

func usage() {
	o := new(tabwriter.Writer)
	o.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s SUBCOMMAND KEYFILE\n\n", os.Args[0])
	fmt.Fprintln(o, "Subcommands:\n")
	fmt.Fprintln(o, "help\tthis help")
	fmt.Fprintln(o, "keygen\tcreate a new keyfile to put into KEYFILE")
	fmt.Fprintln(o, "init\tprepare git repo to use KEYFILE")
	fmt.Fprintln(o, "phrase\tprompts for a phrase to turn into a key")
	fmt.Fprintln(o, "check\tchecks validity of key")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "Plumbing commands (probably don't use):\n")
	fmt.Fprintln(o, "clean\tgit clean filter")
	fmt.Fprintln(o, "smudge\tgit smudge filter")
	fmt.Fprintln(o, "diff\tdiff encrypted files")
	fmt.Fprintln(o)
	o.Flush()
}

func main() {
	var err error
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}
	exe, err = filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to make sense of path: %v\n", err)
		os.Exit(1)
	}
	keyfile, err = filepath.Abs(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to make sense of path: %v\n", err)
		os.Exit(1)
	}
	switch os.Args[1] {
	case "keygen":
		err = keygen()
	case "check":
		err = checkKey()
	case "init":
		err = initRepo()
	case "phrase":
		err = keygenFromPhrase()
	case "clean":
		err = clean()
	case "smudge":
		err = smudge()
	case "diff":
		err = diff(os.Args[3])
	default:
		usage()
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error!: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func keygen() error {
	k, err := NewKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("failure generating key: %v", err)
	}
	return WriteKey(keyfile, k)
}

func keygenFromPhrase() error {
	hkdf := hkdf.New(sha256.New, readPhrase(), nil, nil)
	k, err := NewKey(hkdf)
	if err != nil {
		return fmt.Errorf("failure generating key: %v", err)
	}
	return WriteKey(keyfile, k)
}

func initRepo() error {
	cfgs := [][]string{
		[]string{"git", "config", "filter.grypt.smudge", fmt.Sprintf("%s smudge %s", exe, keyfile)},
		[]string{"git", "config", "filter.grypt.clean", fmt.Sprintf("%s clean %s", exe, keyfile)},
		[]string{"git", "config", "filter.grypt.textconv", fmt.Sprintf("%s textconv %s", exe, keyfile)},
	}

	// check if we have a HEAD
	r := exec.Command("git", "rev-parse", "HEAD")
	hasHEAD := true
	if r.Run() != nil {
		hasHEAD = false
	}

	// check if repo is dirty
	r = exec.Command("git", "-uno", "--porcelain")
	o := new(bytes.Buffer)
	r.Stdout = o
	if r.Run(); o.Len() != 0 && hasHEAD {
		return fmt.Errorf("working directory not clean, stash or merge changes before running `init'")
	}

	// set config options
	for _, cmd := range cfgs {
		c := exec.Command(cmd[0], cmd[1:]...)
		if c.Run() != nil {
			return fmt.Errorf("unable to set config option: %s", cmd)
		}
	}

	// run a forced checkout to decrypt any encrypted files
	if hasHEAD {
		c := exec.Command("git", "checkout", "-f", "HEAD")
		if c.Run() != nil {
			return fmt.Errorf("`git checkout' failed\n%s", c.Stderr)
		}
	}

	fmt.Println(attributesHelp)
	return nil
}

func checkKey() error {
	_, err := ReadKey(keyfile)
	if err != nil {
		os.Exit(1)
	}
	return nil
}

func clean() error {
	k, err := ReadKey(keyfile)
	if err != nil {
		return fmt.Errorf("error reading key: %v", err)
	}
	return Encrypt(os.Stdin, os.Stdout, k)
}

func smudge() error {
	k, err := ReadKey(keyfile)
	if err != nil {
		return fmt.Errorf("error reading key: %v", err)
	}
	return Decrypt(os.Stdin, os.Stdout, k)
}

func diff(f string) error {
	k, err := ReadKey(keyfile)
	if err != nil {
		return err
	}
	file, err := os.Open(f)
	if err != nil {
		return err
	}
	defer file.Close()
	return Decrypt(file, os.Stdout, k)
}
