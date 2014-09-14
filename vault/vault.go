package vault

import (
	"bytes"
	"crypto/cipher"
	"crypto/hmac"
	"fmt"
	"io"

	grypt "polydawn.net/grypt"
)

type Content struct {
	Scheme     grypt.Scheme
	ciphertext []byte // buffer might turn out more appropriate
	IV         []byte // probably these should just go as a part of the b64'd ciphertext?
	MAC        []byte // probably these should just go as a part of the b64'd ciphertext?
}

func (c Content) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "----- BEGIN GRYPT CIPHERTEXT HEADER -----\n")
	fmt.Fprintf(&buf, "Scheme: %s\n", c.Scheme)
	fmt.Fprintf(&buf, "----- END GRYPT CIPHERTEXT HEADER -----\n")
	buf.Write(c.ciphertext)
	buf.WriteRune('\n')
	buf.WriteRune('\n')
	return buf.Bytes(), nil
}

func (c *Content) UnmarshalBinary(data []byte) error {
	return nil
}

// didn't actually intend to ripe this much out of the main package yet, but i'm having import cycle butthurtz, so here we go
func Encrypt(i io.Reader, o io.Writer, k grypt.Key) error {
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
	// header, err := asn1.Marshal(Header{k.Scheme, iv, hmacMsg.Sum(nil)})
	header, err := Content{k.Scheme, iv, hmacMsg.Sum(nil), ciphertext.Bytes()}.MarshalBinary()
	_, err = o.Write(header)
	// if err != nil {
	// 	return err
	// }
	// _, err = io.Copy(o, ciphertext)
	return err
}
