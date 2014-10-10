package schema

import (
	"crypto/cipher"
	"io"

	"code.google.com/p/go.crypto/blowfish"
	"polydawn.net/grypt/ext/blake2b"
)

/*
	@implements Schema
*/
type Blowfish448blake2512ctr struct{}

func (s Blowfish448blake2512ctr) Name() string   { return "blowfish448blake2512ctr" }
func (s Blowfish448blake2512ctr) KeySize() int   { return 56 }
func (s Blowfish448blake2512ctr) MACSize() int   { return 64 }
func (s Blowfish448blake2512ctr) BlockSize() int { return blowfish.BlockSize }

func (s Blowfish448blake2512ctr) NewKey(entropy io.Reader) (Key, error) {
	return newKey(s, entropy)
}

func (s Blowfish448blake2512ctr) Encrypt(input io.Reader, output io.Writer, k Key) error {
	return buildEncrypter(
		s,
		blake2b.New512,
		blowfish.NewCipher,
		cipher.NewCTR,
	)(input, output, k)
}

func (s Blowfish448blake2512ctr) Decrypt(input io.Reader, output io.Writer, k Key) error {
	// CTR mode is an interesting degenerate case because it's literally the same stream for encryption and decryption: the stream is just XOR'd.
	return buildDecrypter(
		s,
		blake2b.New512,
		blowfish.NewCipher,
		cipher.NewCTR,
	)(input, output, k)
}
