package vault

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"io"

	"polydawn.net/grypt/schema"
)

const grypt_vault_typestring = "GRYPT CIPHERTEXT HEADER"

/*
	Certain keys are well-known.  Those are declared in consts in this package.
	Other keys will be passed through when grypt decrypts, alters cleartext payload, and recrypts, even if grypt does not recognize them as well-known keys.

	Headers prefixed with "Grypt-" are reserved for -- you guessed it -- grypt.  They should not be used by any other extensions.
*/
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

type Headers map[string]string

/*
	in: a stream expected to contain grypt headers (in PEM format) followed by the ciphertext binary.
	out: a stream of the cleartext (the headers are returned).
*/
func OpenEnvelope(in io.Reader, out io.Writer, k schema.Key) Headers {
	// read headers
	// note this involves some really shitty and arbitrary assumptions -- like, your header won't be longer than a meg and nobody cares if we overread on the input -- because this pem implementation doesn't have a streamable reader.
	// i understand this interface probably predates a lot of other more refined parts of the standard library, and i understand the no-breaking-changees desires at this point, but boy does it hurt to look at it, and all the more given how wonderfully well thought out the other encoding interfaces are.
	shitbuf := make([]byte, 1024*1024)
	n, err := in.Read(shitbuf)
	if err != nil && err != io.EOF {
		panic(err)
	}
	headerBlock, rest := pem.Decode(shitbuf[:n])
	// OH MY GOD this is a deserialization function that has no error returns.  could this possibly be more wrong.
	if headerBlock == nil {
		panic(fmt.Errorf("not a valid grypt ciphertext file -- no headers detectable"))
	}

	// check version, parse scheme, validate headers in general
	var schemeName string
	var ok bool // ... i don't actually want to keep this, but i can't use `:= later or i'll overshadow the *other* return that i do want to keep.  -.-
	if _, ok := headerBlock.Headers[Header_grypt_version]; !ok {
		// we can switch on this in the future, but right now, we don't know versioned behaviors
		panic(fmt.Errorf("not a valid grypt ciphertext file -- missing required header \"%s\"", Header_grypt_version))
	}
	if schemeName, ok = headerBlock.Headers[Header_grypt_scheme]; !ok {
		panic(fmt.Errorf("not a valid grypt ciphertext file -- missing required header \"%s\"", Header_grypt_scheme))
	}
	if _, ok := headerBlock.Headers[Header_grypt_keyring]; !ok {
		// TODO: keyrings
		panic(fmt.Errorf("not a valid grypt ciphertext file -- missing required header \"%s\"", Header_grypt_keyring))
	}
	sch := schema.ParseSchema(schemeName)

	// push the remainder of body through the cipher
	sch.Decrypt(io.MultiReader(bytes.NewBuffer(rest), in), out, k)

	return headerBlock.Headers
}

/*
	in: a stream of cleartext.
	out: a stream of the headers (in PEM format) followed by the ciphertext binary.
*/
func SealEnvelope(in io.Reader, out io.Writer, k schema.Key) {
	// assemble and output header
	// TODO: support for extra headers, we currently fail at passthrough of headers we don't recognize and that's amateur horseshit
	headerBlock := &pem.Block{
		Type: grypt_vault_typestring,
		Headers: Headers{
			Header_grypt_version: "1.0",
			Header_grypt_scheme:  k.Scheme.Name(),
			Header_grypt_keyring: "default", // :/ TODO: keyrings
		},
		// pem.Block.Bytes is a zero value for us, we're not gonna use b64
	}
	headerBytes := pem.EncodeToMemory(headerBlock)
	if _, err := out.Write(headerBytes); err != nil {
		panic(err)
	}

	// push the remainder of body through the cipher
	if err := k.Scheme.Encrypt(in, out, k); err != nil {
		panic(err)
	}
}
