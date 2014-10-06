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
type Aes256sha256ctr struct {}

func (s Aes256sha256ctr) KeySize() int {
	return 32
}

func (s Aes256sha256ctr) MACSize() int {
	return 32
}

func (s Aes256sha256ctr) NewKey(entropy io.Reader) (Key, error) {
	var err error
	symKey := make([]byte, s.KeySize())
	macKey := make([]byte, s.MACSize())
	_, err = io.ReadFull(entropy, macKey)
	if err != nil {
		return Key{}, err
	}
	_, err = io.ReadFull(entropy, symKey)
	if err != nil {
		return Key{}, err
	}
	return Key{symKey, macKey}, nil
}

/*
	Output is:
	- IV (fixed length)
	- length of ciphertext body (fixed length, 8 byte signed long)
	- ciphertext body
	- MAC of ciphertext body (fixed length)

	Other headers like which schema this is, etc, are expected to be kept elsewhere as necessary.
*/
func (s Aes256sha256ctr) Encrypt(input io.Reader, output io.Writer, k Key) error {
	return encrypt(s, sha256.New, aes.NewCipher, cipher.NewCTR)(input, output, k)
}
