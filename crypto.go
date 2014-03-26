package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"encoding/asn1"
	"fmt"
	"io"
)

// Decrypt ciphertext into plaintext.
func Decrypt(i io.Reader, o io.Writer, k Key) error {
	header := Header{}
	headBuf := new(bytes.Buffer)
	c, err := aes.NewCipher(k.Key)
	if err != nil {
		return fmt.Errorf("unabled to create aes cipher: %v", err)
	}

	// Read a small chunk and try to parse the header
	_, err = io.CopyN(headBuf, i, 1024)
	if err != nil && err != io.EOF {
		return fmt.Errorf("copyN error: %v", err)
	}
	rest, err := asn1.Unmarshal(headBuf.Bytes(), &header)
	if err != nil {
		return fmt.Errorf("unmarshal error: %v", err)
	}
	if header.Scheme != k.Scheme {
		return fmt.Errorf("key is unable to decrypt this data")
	}
	s := cipher.StreamWriter{
		S: cipher.NewCTR(c, header.MAC[:k.Scheme.BlockSize()]),
		W: o,
	}

	// read the encrypted file and decrypt
	_, err = io.Copy(s, io.MultiReader(bytes.NewBuffer(rest), i))
	return err
}

// Encrypt plaintext to ciphertext.
func Encrypt(i io.Reader, o io.Writer, k Key) error {
	buf := new(bytes.Buffer)
	c, err := aes.NewCipher(k.Key)
	if err != nil {
		return err
	}
	h := hmac.New(k.Scheme.Hash(), k.HMAC)
	w := io.MultiWriter(buf, h)

	// Read in the file, calculating the hmac and buffering it
	if _, err := io.Copy(w, i); err != nil {
		return err
	}
	mac := h.Sum(nil)

	s := cipher.StreamWriter{
		S: cipher.NewCTR(c, mac[:k.Scheme.BlockSize()]),
		W: o,
	}

	// serialize our header and append the encrypted file
	header, err := asn1.Marshal(Header{k.Scheme, mac})
	if err != nil {
		return err
	}
	_, err = o.Write(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(s, buf)
	return err
}
