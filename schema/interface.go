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
	// every schema can be looked up by its own name (obviously), and these are the ones that end up serialized in headers
	for _, s := range []Schema{
		Aes256sha256ctr{},
		Aes256keccak256ctr{},
		Aes256blake2256ctr{},
		Blowfish448sha256ctr{},
		Blowfish448blake2512ctr{},
	} {
		schemas[s.Name()] = s
	}
	// additional names map onto the some things.
	schemas["default"] = schemas[Aes256sha256ctr{}.Name()]
	schemas["keccak"] = schemas[Aes256keccak256ctr{}.Name()]
	schemas["blake2"] = schemas[Aes256blake2256ctr{}.Name()]
	schemas["blowfish"] = schemas[Blowfish448sha256ctr{}.Name()]
	schemas["blakefish"] = schemas[Blowfish448blake2512ctr{}.Name()]
}

func ParseSchema(s string) Schema {
	if s, ok := schemas[s]; ok {
		return s
	} else {
		panic(fmt.Errorf("invalid encryption schema name"))
	}
}
