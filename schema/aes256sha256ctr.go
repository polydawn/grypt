package schema

import (
	"bytes"
	"crypto/cipher"
	"crypto/hmac"
	"io"

	grypt "polydawn.net/grypt"
)

type Aes256sha256ctr struct {}

func (s Aes256sha256ctr) KeySize() int {
	return 32
}

func (s Aes256sha256ctr) MACSize() int {
	return 32
}

func (s Aes256sha256ctr) NewKey(entropy io.Reader) (grypt.Key, error) {
	var err error
	symKey := make([]byte, s.KeySize())
	macKey := make([]byte, s.MACSize())
	_, err = io.ReadFull(entropy, macKey)
	if err != nil {
		return grypt.Key{}, err
	}
	_, err = io.ReadFull(entropy, symKey)
	if err != nil {
		return grypt.Key{}, err
	}
	return grypt.Key{grypt.Scheme(-1), symKey, macKey}, nil //FIXME: that -1 in the struct is pants, we just don't need that in the key struct anymore this way
}

/*
	Output is:
	- IV (fixed length)
	- length of ciphertext body (fixed length, 8 byte signed long)
	- ciphertext body
	- MAC of ciphertext body (fixed length)

	Other headers like which schema this is, etc, are expected to be kept elsewhere as necessary.
*/
func (s Aes256sha256ctr) Encrypt(input io.Reader, output io.Writer, k grypt.Key) error {
	// Read in the file, calculating the IV and buffering it
	// impl note: this won't do well with gig files... but this could easily be replaced with re-reading the input, if we had knew we had a resettable reader like a disk
	plaintext := new(bytes.Buffer)
	hmacIV := hmac.New(k.Scheme.Hash(), k.HMAC)
	mw := io.MultiWriter(plaintext, hmacIV)

	if _, err := io.Copy(mw, input); err != nil {
		return err
	}
	iv := hmacIV.Sum(nil)[:k.Scheme.BlockSize()]

	// commit the iv to output
	if _, err := output.Write(iv); err != nil {
		return err
	}

	// initialize cipher, hmac, and write the expected ciphertext size
	hmacMsg := hmac.New(k.Scheme.Hash(), k.HMAC)
	blockCipher, err := k.Scheme.NewCipher(k.Key)
	if err != nil {
		return err
	}
	//TODO the ciphertext size

	// push the input body through the cipher, and the ciphertext both out and through the hmac
	mw = io.MultiWriter(output, hmacMsg)
	streamCipher := cipher.StreamWriter{
		S: cipher.NewCTR(blockCipher, iv),
		W: mw,
	}
	_, err = io.Copy(streamCipher, plaintext)
	if err != nil {
		return err
	}
	// TODO: verify: the stream cipher had better damn well know how to do padding on close

	// ciphertext body now all pushed; commit the hmac to output and we're done
	if _, err := output.Write(hmacMsg.Sum(nil)); err != nil {
		return err
	}
	return nil
}
