package schema

import (
	"crypto/aes"
	"crypto/cipher"
	"io"

	"polydawn.net/grypt/ext/blake2b"
)

var _ Schema = Aes256blake2256ctr{}

type Aes256blake2256ctr struct{}

func (s Aes256blake2256ctr) Name() string   { return "aes256blake2256ctr" }
func (s Aes256blake2256ctr) KeySize() int   { return 32 }
func (s Aes256blake2256ctr) MACSize() int   { return 32 }
func (s Aes256blake2256ctr) BlockSize() int { return aes.BlockSize }

func (s Aes256blake2256ctr) NewKey(entropy io.Reader) (Key, error) {
	return newKey(s, entropy)
}

func (s Aes256blake2256ctr) Encrypt(input io.Reader, output io.Writer, k Key) error {
	return buildEncrypter(
		s,
		blake2b.New256,
		aes.NewCipher,
		cipher.NewCTR,
	)(input, output, k)
}

func (s Aes256blake2256ctr) Decrypt(input io.Reader, output io.Writer, k Key) error {
	// CTR mode is an interesting degenerate case because it's literally the same stream for encryption and decryption: the stream is just XOR'd.
	return buildDecrypter(
		s,
		blake2b.New256,
		aes.NewCipher,
		cipher.NewCTR,
	)(input, output, k)
}
