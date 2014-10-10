package schema

import "io"

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
