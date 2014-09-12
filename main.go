package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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
	encryptionScheme Scheme
	schemeString     = flag.String("t", "default", "Which encryption scheme to use (only applicable to 'phrase' and 'keygen')")
)

type (
	// Header contains information about the encryption scheme
	Header struct {
		Scheme Scheme
		IV     []byte
		MAC    []byte
	}
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s [OPTIONS] SUBCOMMAND KEYFILE\n", os.Args[0])
	fmt.Fprintln(os.Stderr, `
SUBCOMMANDS:

help    this help
keygen  create a new keyfile to put into KEYFILE
init    prepare git repo to use KEYFILE
phrase  prompts for a phrase to turn into a key
check   checks validity of key

OPTIONS:
`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, `
Valid encryption schemes are:

 * AES-256/SHA-256          (default, aes256sha256)
 * AES-256/Keccak-256       (keccak, aes256keccak256)
 * AES-256/BLAKE2-256       (blake2, aes256blake2256)
 * Blowfish-448/SHA-256     (blowfish, blowfish448sha256)
 * Blowfish-448/BLAKE2-512  (blakefish, blowfish448blake2512)
`)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	var err error
	if len(flag.Args()) < 1 {
		usage()
		os.Exit(1)
	}
	exe = os.Args[0]
	keyfile, err = filepath.Abs(flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to make sense of path: %v\n", err)
		os.Exit(1)
	}
	encryptionScheme, err = ParseScheme(*schemeString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine encryption scheme: %v", err)
		os.Exit(2)
	}
	switch flag.Arg(0) {
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
		err = diff(flag.Arg(2))
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
	k, err := NewKey(rand.Reader, encryptionScheme)
	if err != nil {
		return fmt.Errorf("failure generating key: %v", err)
	}
	return WriteKey(keyfile, k)
}

func keygenFromPhrase() error {
	hkdf := hkdf.New(sha256.New, readPhrase(), nil, nil)
	k, err := NewKey(hkdf, encryptionScheme)
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
	r = exec.Command("git", "status", "-uno", "--porcelain")
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
