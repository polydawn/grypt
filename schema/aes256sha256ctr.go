package schema

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"io"
)

/*
	@implements Schema
*/
type Aes256sha256ctr struct{}

func (s Aes256sha256ctr) Name() string   { return "aes256sha256ctr" }
func (s Aes256sha256ctr) KeySize() int   { return 32 }
func (s Aes256sha256ctr) MACSize() int   { return 32 }
func (s Aes256sha256ctr) BlockSize() int { return aes.BlockSize }

func (s Aes256sha256ctr) NewKey(entropy io.Reader) (Key, error) {
	return newKey(s, entropy)
}

func (s Aes256sha256ctr) Encrypt(input io.Reader, output io.Writer, k Key) error {
	return buildEncrypter(
		s,
		sha256.New,
		aes.NewCipher,
		cipher.NewCTR,
	)(input, output, k)
}

func (s Aes256sha256ctr) Decrypt(input io.Reader, output io.Writer, k Key) error {
	// CTR mode is an interesting degenerate case because it's literally the same stream for encryption and decryption: the stream is just XOR'd.
	return buildDecrypter(
		s,
		sha256.New,
		aes.NewCipher,
		cipher.NewCTR,
	)(input, output, k)
}
