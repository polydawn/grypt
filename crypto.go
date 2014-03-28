package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/hmac"
	"encoding/asn1"
	"fmt"
	"io"
)

/*
The files we write/read have a small header tacked on (see type Header and
key.go) that carries some encryption scheme information and relevant nonces.

The header format should probably stop changing once We are happy with how the
encryption works.
*/

// Decrypt ciphertext into plaintext.
func Decrypt(i io.Reader, o io.Writer, k Key) error {
	header := Header{}
	headBuf := new(bytes.Buffer)
	decBuf := new(bytes.Buffer)
	h := hmac.New(k.Scheme.Hash(), k.HMAC)
	c, err := k.Scheme.NewCipher(k.Key)
	if err != nil {
		return fmt.Errorf("unabled to create cipher: %v", err)
	}

	// Read a small chunk and try to parse the header
	// This is an arbitrary number.
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
		S: cipher.NewCTR(c, header.IV),
		W: decBuf,
	}
	mw := io.MultiWriter(s, h)

	// read the encrypted file and decrypt
	_, err = io.Copy(mw, io.MultiReader(bytes.NewBuffer(rest), i))
	if err != nil {
		return err
	}
	if !hmac.Equal(header.MAC, h.Sum(nil)) {
		return fmt.Errorf("unable to verify file")
	}
	_, err = io.Copy(o, decBuf)
	return err
}

// Encrypt plaintext to ciphertext.
//
// TODO: meditate on improvements. There seems like one too many buffers.
func Encrypt(i io.Reader, o io.Writer, k Key) error {
	plaintext := new(bytes.Buffer)
	ciphertext := new(bytes.Buffer)
	c, err := k.Scheme.NewCipher(k.Key)
	if err != nil {
		return err
	}
	hmacIV := hmac.New(k.Scheme.Hash(), k.HMAC)
	hmacMsg := hmac.New(k.Scheme.Hash(), k.HMAC)
	mw := io.MultiWriter(plaintext, hmacIV)

	// Read in the file, calculating the IV and buffering it
	if _, err := io.Copy(mw, i); err != nil {
		return err
	}
	iv := hmacIV.Sum(nil)[:k.Scheme.BlockSize()]

	// write ciphertext into buffer and the hmac
	mw = io.MultiWriter(ciphertext, hmacMsg)
	s := cipher.StreamWriter{
		S: cipher.NewCTR(c, iv),
		W: mw,
	}
	_, err = io.Copy(s, plaintext)
	if err != nil {
		return err
	}

	// serialize our header and append the encrypted file
	header, err := asn1.Marshal(Header{k.Scheme, iv, hmacMsg.Sum(nil)})
	_, err = o.Write(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(o, ciphertext)
	return err
}
