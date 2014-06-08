package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"os"

	"code.google.com/p/go.crypto/blowfish"
	"code.google.com/p/go.crypto/sha3"
	"polydawn.net/grypt/ext/blake2b"
)

const (
	// Use AES-256 with a SHA-256 HMAC
	AES256_SHA256 Scheme = iota
	// Use AES-256 with a Keccak-256 (SHA3) HMAC
	AES256_Keccak256
	// Use Blowfish-448 with a SHA-256 HMAC
	Blowfish448_SHA256
	// Use AES-265 with a BLAKE2-256 HMAC
	AES256_BLAKE2256
	// Use Blowfish-448 with a BLAKE2-512 HMAC
	Blowfish448_BLAKE2512
)

var (
	// Scheme used for new keys.
	DefaultScheme = AES256_SHA256
	// The scheme indicated does not exist or is not supported.
	ErrInvalidScheme = fmt.Errorf("invalid scheme")
)

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

func ParseScheme(s string) (Scheme, error) {
	switch s {
	case "default", "aes256sha256":
		return AES256_SHA256, nil
	case "keccak", "aes256keccak256":
		return AES256_Keccak256, nil
	case "blowfish", "blowfish448sha256":
		return Blowfish448_SHA256, nil
	case "blake2", "aes256blake2256":
		return AES256_BLAKE2256, nil
	case "blakefish", "blowfish448blake2512":
		return Blowfish448_BLAKE2512, nil
	}
	return Scheme(-1), ErrInvalidScheme
}

func (s Scheme) KeySize() int {
	switch s {
	case AES256_SHA256, AES256_Keccak256, AES256_BLAKE2256:
		return 32
	case Blowfish448_SHA256, Blowfish448_BLAKE2512:
		return 56
	default:
		panic("invalid Scheme")
	}
}

func (s Scheme) MACSize() int {
	switch s {
	case Blowfish448_SHA256, AES256_SHA256, AES256_Keccak256, AES256_BLAKE2256:
		return 32
	case Blowfish448_BLAKE2512:
		return 64
	default:
		panic("invalid Scheme")
	}
}

func (s Scheme) BlockSize() int {
	switch s {
	case AES256_SHA256, AES256_Keccak256, AES256_BLAKE2256:
		return aes.BlockSize
	case Blowfish448_SHA256, Blowfish448_BLAKE2512:
		return blowfish.BlockSize
	default:
		panic("invalid Scheme")
	}
}

// Returns a cipher.Block of the relevant cipher
func (s Scheme) NewCipher(key []byte) (cipher.Block, error) {
	switch s {
	case AES256_SHA256, AES256_Keccak256, AES256_BLAKE2256:
		return aes.NewCipher(key)
	case Blowfish448_SHA256, Blowfish448_BLAKE2512:
		return blowfish.NewCipher(key)
	default:
		panic("invalid Scheme")
	}
}

// Returns '.New' of the relevant hash package
func (s Scheme) Hash() func() hash.Hash {
	switch s {
	case Blowfish448_SHA256, AES256_SHA256:
		return sha256.New
	case AES256_Keccak256:
		return sha3.NewKeccak256
	case AES256_BLAKE2256:
		return blake2b.New256
	case Blowfish448_BLAKE2512:
		return blake2b.New512
	default:
		panic("invalid Scheme")
	}
}

func (s Scheme) String() string {
	switch s {
	case Blowfish448_SHA256:
		return "Blowfish-448/SHA-256"
	case AES256_SHA256:
		return "AES-256/SHA-256"
	case AES256_Keccak256:
		return "AES-256/Keccak-256"
	case AES256_BLAKE2256:
		return "AES-256/BLAKE2-256"
	case Blowfish448_BLAKE2512:
		return "Blowfish-448/BLAKE2-512"
	default:
		panic("invalid Scheme")
	}
}

// return a Key from the supplied Reader
func NewKey(r io.Reader, s Scheme) (Key, error) {
	var err error
	symKey := make([]byte, s.KeySize())
	macKey := make([]byte, s.MACSize())
	_, err = io.ReadFull(r, macKey)
	if err != nil {
		return Key{}, err
	}
	_, err = io.ReadFull(r, symKey)
	if err != nil {
		return Key{}, err
	}
	return Key{s, symKey, macKey}, nil
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
