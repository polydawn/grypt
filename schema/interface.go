package schema

import "io"

type Schema interface {
	Name() string

	KeySize() int
	MACSize() int

	NewKey(entropy io.Reader) (Key, error)

	Encrypt(input io.Reader, output io.Writer, k Key) error
	Decrypt(input io.Reader, output io.Writer, k Key) error
}
