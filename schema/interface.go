package schema

import (
	"io"
	"fmt"
)

type Schema interface {
	Name() string

	KeySize() int
	MACSize() int
	BlockSize() int

	NewKey(entropy io.Reader) (Key, error)

	Encrypt(input io.Reader, output io.Writer, k Key) error
	Decrypt(input io.Reader, output io.Writer, k Key) error
}

/*
	Key struct stores the two byte slices for most symmetric crypto operations:
	the cipher key and the hmac key.

	This is a simplifying assumption for all the interfaces we currently use, but may break
	for other kinds of (very) exotic cipher suites we don't yet support.
*/
type Key struct {
	cipherKey []byte
	hmacKey   []byte
}

var schemas map[string]Schema = make(map[string]Schema)

func init() {
	for _, s := range []Schema{
		Aes256sha256ctr{},
		Aes256keccak256ctr{},
		Blowfish448sha256ctr{},
		Aes256blake2256ctr{},
		Blowfish448blake2512ctr{},
	} {
		schemas[s.Name()] = s
	}
}

func ParseSchema(s string) Schema {
	if s, ok := schemas[s]; ok {
		return s
	} else {
		panic(fmt.Errorf("invalid encryption schema name"))
	}
}
