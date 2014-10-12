package schema

import (
	"fmt"
	"io"
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

var schemas map[string]Schema = make(map[string]Schema)
var extraNames map[string]Schema = make(map[string]Schema)

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
	extraNames["default"] = schemas[Aes256sha256ctr{}.Name()]
	extraNames["keccak"] = schemas[Aes256keccak256ctr{}.Name()]
	extraNames["blake2"] = schemas[Aes256blake2256ctr{}.Name()]
	extraNames["blowfish"] = schemas[Blowfish448sha256ctr{}.Name()]
	extraNames["blakefish"] = schemas[Blowfish448blake2512ctr{}.Name()]
}

func ParseSchema(s string) Schema {
	if sch, ok := schemas[s]; ok {
		return sch
	} else if sch, ok := extraNames[s]; ok {
		return sch
	} else {
		panic(fmt.Errorf("invalid encryption schema name"))
	}
}
