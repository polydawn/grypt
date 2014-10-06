package schema

import (
	"bytes"
	"crypto/cipher"
	"crypto/hmac"
	"hash"
	"io"
)

type encrypter func(input io.Reader, output io.Writer, k Key) error

/*
	Generic symmetric encryption for:
	  - a deterministic content-based IV
	  - symmetric ciphertext stream (derived from a block cipher via a mode of operation)
	  - MAC over ciphertext
	Choice of block cipher, MAC construction, and stream mode are parameters.

	Output is:
	- IV (fixed length)
	- length of ciphertext body (fixed length, 8 byte signed long)
	- ciphertext body
	- MAC of ciphertext body (fixed length)

	Other headers like which schema this is, etc, are expected to be kept elsewhere as necessary.
*/
func encrypt(sch Schema, macFactory func() hash.Hash, cipherFactory func(key []byte) (cipher.Block, error), streamFactory func (cipher.Block, []byte) cipher.Stream) encrypter {
	// implementation note: the golang stdlib distinction between cipher.Stream and cipher.BlockMode is... odd, and could be readily wallpapered over with some really derpy wrappers.  haven't bothered yet.

	// also: feel slightly bad about passing in both the schema, and all its functors.  could make the Schema interface also return all these things.
	// really wish the hash/block/stream classes just returned their own sizes consistently, because that would make this a nonissue.

	// @implements encrypter
	return func(input io.Reader, output io.Writer, k Key) error {
		// Read in the file, calculating the IV and buffering it
		// impl note: this won't do well with gig files... but this could easily be replaced with re-reading the input, if we had knew we had a resettable reader like a disk
		plaintext := new(bytes.Buffer)
		hmacIV := hmac.New(macFactory, k.hmacKey)
		mw := io.MultiWriter(plaintext, hmacIV)

		if _, err := io.Copy(mw, input); err != nil {
			return err
		}
		iv := hmacIV.Sum(nil)[:sch.MACSize()]

		// commit the iv to output
		if _, err := output.Write(iv); err != nil {
			return err
		}

		// initialize cipher, hmac, and write the expected ciphertext size
		hmacMsg := hmac.New(macFactory, k.hmacKey)
		blockCipher, err := cipherFactory(k.cipherKey)
		if err != nil {
			return err
		}
		//TODO the ciphertext size

		// push the input body through the cipher, and the ciphertext both out and through the hmac
		mw = io.MultiWriter(output, hmacMsg)
		streamCipher := cipher.StreamWriter{
			S: streamFactory(blockCipher, iv),
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
}
