package vault

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"crypto/hmac"
	"fmt"
	"io"
	"regexp"
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

var Nonoptional_headers = []string{
	Header_grypt_version,
	Header_grypt_scheme,
	Header_grypt_keyring,
}

var header_regexp = regexp.MustCompile("^([[:upper:]][[:alnum:]-]*):[[:space:]]*([[:alnum:][:punct:]]*)[[:space:]]*")

type Content struct {
	Headers    Headers
	ciphertext []byte // buffer might turn out more appropriate
	IV         []byte // probably these should just go as a part of the b64'd ciphertext?
	MAC        []byte // probably these should just go as a part of the b64'd ciphertext?
}

const (
	grypt_vault_start = "----- BEGIN GRYPT CIPHERTEXT HEADER -----\n"
	grypt_vault_end   = "----- END GRYPT CIPHERTEXT HEADER -----\n"
)

func (c Content) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer

	// start header
	fmt.Fprintf(&buf, grypt_vault_start)

	// place certain non-optional headers at the front of the parade.
	for _, key := range Nonoptional_headers {
		value, ok := c.Headers[key]
		if !ok {
			return nil, fmt.Errorf("vault requires header %s", key)
		}
		fmt.Fprintf(&buf, "%s: %s\n", key, value)
	}

	// iterate over remaining headers.  skip the ones already included above.
L1:
	for key, value := range c.Headers {
		for _, already := range Nonoptional_headers {
			if key == already {
				continue L1
			}
		}
		fmt.Fprintf(&buf, "%s: %s\n", key, value)
	}

	// end header
	fmt.Fprintf(&buf, grypt_vault_end)

	// drop ciphertext.  length is embedded in the binary form.
	buf.Write(c.ciphertext)

	// trailing whitespace to moderately decrease the odds of your terminal crying if you cat this file.
	buf.WriteRune('\n')
	buf.WriteRune('\n')

	return buf.Bytes(), nil
}

func (c *Content) UnmarshalBinary(data []byte) error {
	reader := bufio.NewReader(bytes.NewBuffer(data))
	var line string
	var err error

	// first line absolutely must be our header
	if line, err = reader.ReadString('\n'); err != nil {
		return err
	}
	if line != grypt_vault_start {
		return fmt.Errorf("invalid grypt vault header: doesn't look like grypt vault ciphertext")
	}

	// read one line at a time until we see the end of our header
	var rows []string
	for {
		if line, err = reader.ReadString('\n'); err != nil {
			return err
		}
		rows = append(rows, line)
		if line == grypt_vault_end {
			break
		}
	}

	// parse all those header entries
	c.Headers = Headers{}
	for _, row := range rows {
		matches := header_regexp.FindStringSubmatch(row)
		if len(matches) != 3 {
			continue // this isn't one of ours
			// consider adding a format verification step as a CheckKey subcommand that could warn about bananas formats, in case someone starts writing their own headers for some unearthly reason
		}
		c.Headers[matches[1]] = matches[2]
	}

	// we can now act like a reader for the payload decryption (which should know enough about its own format to not over-read; we may have more bytes trailing than just payload)

	return nil
}

// didn't actually intend to ripe this much out of the main package yet, but i'm having import cycle butthurtz, so here we go
func Encrypt(ctx grypt.Context, i io.Reader, o io.Writer, k grypt.Key) error {
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
		Header_grypt_version: ctx.GryptVersion,
		Header_grypt_scheme:  fmt.Sprintf("%s", k.Scheme), // TODO this does roughly "what I mean", but should probably be replaced by a marshaller spec on a solid scheme type
		Header_grypt_keyring: ctx.Keyring,
	}
	serial, err := Content{headers, iv, hmacMsg.Sum(nil), ciphertext.Bytes()}.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = o.Write(serial)
	return err
}
