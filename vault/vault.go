package vault

import (
	"bytes"
	"crypto/cipher"
	"crypto/hmac"
	"fmt"
	"io"

	grypt "polydawn.net/grypt"
)

/*
	Grypt content vault format headers are simple string key-value pairs.
	They roughly follow the outline of git commit message suffixes:

	- keys should be formatted as "Title-Case-Phrases:" terminated with a colon, then whitespace.
	  - keys may not contain whitespace or colon characters.
	  - keys are case sensitive.
	  - it is strongly recommended that keys only contain alphanum and dashes.
	- values follow keys, and must be a single line.

	Keys and values are both interpretted as UTF-8 sequences.  There is no strict length on either keys or values.

	Certain keys are well-known.  Those are declared in consts in this package.
	Other keys will be passed through when grypt decrypts, alters cleartext payload, and recrypts, even if grypt does not recognize them as well-known keys.
*/
type Headers map[string]string

// ^ note: goddamnit, ordered maps again.  it would certainly be charming to put well-known keys first, and maintain orders of others so diffs aren't randomly shitty.  as it is we'll have to sort and call it a day.

const (
	/*
		Names the version of grypt that produced the current content.
	*/
	Header_grypt_version = "Grypt-Version"

	/*
		Names the encryption and verification scheme used in the ciphertext payload.
		This header is directly necessary for decryption and verification to work correctly.  Do not change its value.
	*/
	Header_grypt_scheme = "Grypt-Scheme"

	/*
		Names the key used to encrypt the ciphertext payload.

		This is not strictly checked (different users may follow different key naming conventions if they so desire),
		but by default warnings will be emitted if decryption fails and the keyname does not match the name of the provided key.
		These warnings are provided to help defuse "D'oh" moments with key management and argument fat-fingering, in the common case where all users of a repo have agreed on the same key names.
	*/
	Header_grypt_keyring = "Grypt-Keyring"
)

type Content struct {
	Headers    Headers
	ciphertext []byte // buffer might turn out more appropriate
	IV         []byte // probably these should just go as a part of the b64'd ciphertext?
	MAC        []byte // probably these should just go as a part of the b64'd ciphertext?
}

func (c Content) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	// start header
	fmt.Fprintf(&buf, "----- BEGIN GRYPT CIPHERTEXT HEADER -----\n")
	// place certain non-optional headers at the front of the parade.  report errors if they are not assigned.
	// TODO
	// iterate over remaining headers.  skip the ones already included above.
	for key, value := range c.Headers {
		fmt.Fprintf(&buf, "%s: %s\n", key, value)
	}
	// end header
	fmt.Fprintf(&buf, "----- END GRYPT CIPHERTEXT HEADER -----\n")
	// drop ciphertext.  length is embedded in the binary form.
	buf.Write(c.ciphertext)
	// trailing whitespace to moderately decrease the odds of your terminal crying if you cat this file.
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
	headers := Headers{
		Header_grypt_scheme: fmt.Sprintf("%s", k.Scheme), // TODO this does roughly "what I mean", but should probably be replaced by a marshaller spec on a solid scheme type
	}
	serial, err := Content{headers, iv, hmacMsg.Sum(nil), ciphertext.Bytes()}.MarshalBinary()
	_, err = o.Write(serial)
	// if err != nil {
	// 	return err
	// }
	// _, err = io.Copy(o, ciphertext)
	return err
}
