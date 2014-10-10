package schema

import (
	"crypto/cipher"
	"crypto/sha256"
	"io"

	"code.google.com/p/go.crypto/blowfish"
)

/*
	@implements Schema
*/
type Blowfish448sha256ctr struct{}

func (s Blowfish448sha256ctr) Name() string   { return "blowfish448sha256ctr" }
func (s Blowfish448sha256ctr) KeySize() int   { return 56 }
func (s Blowfish448sha256ctr) MACSize() int   { return 32 }
func (s Blowfish448sha256ctr) BlockSize() int { return blowfish.BlockSize }

func (s Blowfish448sha256ctr) NewKey(entropy io.Reader) (Key, error) {
	return newKey(s, entropy)
}

func (s Blowfish448sha256ctr) Encrypt(input io.Reader, output io.Writer, k Key) error {
	return buildEncrypter(
		s,
		sha256.New,
		blowfish.NewCipher,
		cipher.NewCTR,
	)(input, output, k)
}

func (s Blowfish448sha256ctr) Decrypt(input io.Reader, output io.Writer, k Key) error {
	// CTR mode is an interesting degenerate case because it's literally the same stream for encryption and decryption: the stream is just XOR'd.
	return buildDecrypter(
		s,
		sha256.New,
		blowfish.NewCipher,
		cipher.NewCTR,
	)(input, output, k)
}
