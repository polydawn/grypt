package vault

import (
//"polydawn.net/grypt/schema"
)

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
