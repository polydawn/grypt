package schema

import (
	"crypto/aes"
	"crypto/cipher"
	"io"

	"code.google.com/p/go.crypto/sha3"
)

var _ Schema = Aes256keccak256ctr{}

type Aes256keccak256ctr struct{}

func (s Aes256keccak256ctr) Name() string   { return "aes256keccak256ctr" }
func (s Aes256keccak256ctr) KeySize() int   { return 32 }
func (s Aes256keccak256ctr) MACSize() int   { return 32 }
func (s Aes256keccak256ctr) BlockSize() int { return aes.BlockSize }

func (s Aes256keccak256ctr) NewKey(entropy io.Reader) (Key, error) {
	return newKey(s, entropy)
}

func (s Aes256keccak256ctr) Encrypt(input io.Reader, output io.Writer, k Key) error {
	return buildEncrypter(
		s,
		sha3.NewKeccak256,
		aes.NewCipher,
		cipher.NewCTR,
	)(input, output, k)
}

func (s Aes256keccak256ctr) Decrypt(input io.Reader, output io.Writer, k Key) error {
	// CTR mode is an interesting degenerate case because it's literally the same stream for encryption and decryption: the stream is just XOR'd.
	return buildDecrypter(
		s,
		sha3.NewKeccak256,
		aes.NewCipher,
		cipher.NewCTR,
	)(input, output, k)
}
