package main

import (
	"bytes"
	"crypto/aes"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"hash"
	"io"
	"os"
)

const (
	// Use AES-256 with a SHA-256 HMAC
	AES256_SHA256 Scheme = iota
)

// Scheme used for new keys. Currently not user tunable.
var DefaultScheme = AES256_SHA256

type (
	// Key is used for symmetric encryption in/out of the blob store
	Key struct {
		Scheme Scheme
		// Symetric cipher key
		Key []byte
		// Key for HMAC
		HMAC []byte
	}
	// Encryption scheme
	Scheme int
)

func (s Scheme) KeySize() int {
	switch s {
	case AES256_SHA256:
		return 32
	default:
		panic("invalid Scheme")
	}
}

func (s Scheme) MACSize() int {
	switch s {
	case AES256_SHA256:
		return 32
	default:
		panic("invalid Scheme")
	}
}

func (s Scheme) BlockSize() int {
	switch s {
	case AES256_SHA256:
		return aes.BlockSize
	default:
		panic("invalid Scheme")
	}
}

// Returns '.New' of the relevant hash package
func (s Scheme) Hash() func() hash.Hash {
	switch s {
	case AES256_SHA256:
		return sha256.New
	default:
		panic("invalid Scheme")
	}
}

// return a Key from the supplied Reader
func NewKey(r io.Reader) (Key, error) {
	var err error
	symKey := make([]byte, DefaultScheme.KeySize())
	macKey := make([]byte, DefaultScheme.MACSize())
	_, err = io.ReadFull(r, macKey)
	if err != nil {
		return Key{}, err
	}
	_, err = io.ReadFull(r, symKey)
	if err != nil {
		return Key{}, err
	}
	return Key{DefaultScheme, symKey, macKey}, nil
}

// base64 encode and write key 'k' to file 'f'
func WriteKey(f string, k Key) error {
	bits, err := asn1.Marshal(k)
	if err != nil {
		return err
	}
	file, err := os.Create(f)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := base64.NewEncoder(base64.StdEncoding, file)
	_, err = bytes.NewBuffer(bits).WriteTo(enc)
	if err != nil {
		return err
	}
	enc.Close()
	return nil
}

// read and decode a key from file 'f'
func ReadKey(f string) (Key, error) {
	k := Key{}
	bits := new(bytes.Buffer)
	file, err := os.Open(f)
	if err != nil {
		return Key{}, err
	}
	defer file.Close()
	dec := base64.NewDecoder(base64.StdEncoding, file)
	_, err = bits.ReadFrom(dec)
	if err != nil {
		return Key{}, err
	}
	_, err = asn1.Unmarshal(bits.Bytes(), &k)
	if err != nil {
		return Key{}, err
	}
	return k, nil
}
