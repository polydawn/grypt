package schema

import (
	"bytes"
	"crypto/cipher"
	"crypto/hmac"
	"encoding/binary"
	"hash"
	"io"
)

type encrypter func(input io.Reader, output io.Writer, k Key) error

type decrypter func(input io.Reader, output io.Writer, k Key) error

/*
	Generic symmetric encryption for:
	  - a deterministic content-based IV
	  - symmetric ciphertext stream (derived from a block cipher via a mode of operation)
	  - MAC over ciphertext
	Choice of block cipher, MAC construction, and stream mode are parameters.

	Output is:
	- IV (fixed length)
	- length of ciphertext body (fixed length, 8 byte signed big-endian long)
	- ciphertext body
	- MAC of ciphertext body (fixed length)

	Other headers like which schema this is, etc, are expected to be kept elsewhere as necessary.
*/
func buildEncrypter(sch Schema, macFactory func() hash.Hash, cipherFactory func(key []byte) (cipher.Block, error), streamFactory func(cipher.Block, []byte) cipher.Stream) encrypter {
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
		iv := hmacIV.Sum(nil)[:sch.BlockSize()]

		// commit the iv to output
		if _, err := output.Write(iv); err != nil {
			return err
		}

		// initialize cipher, hmac
		hmacMsg := hmac.New(macFactory, k.hmacKey)
		blockCipher, err := cipherFactory(k.cipherKey)
		if err != nil {
			return err
		}

		// commit the expected ciphertext size to output
		blocks := plaintext.Len() / sch.BlockSize()
		if plaintext.Len() % sch.BlockSize() != 0 {
			blocks++
		}
		binary.Write(output, binary.BigEndian, int64(blocks*sch.BlockSize()))

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

		// ciphertext body now all pushed; commit the hmac to output and we're done
		if _, err := output.Write(hmacMsg.Sum(nil)); err != nil {
			return err
		}
		return nil
	}
}

func buildDecrypter(sch Schema, macFactory func() hash.Hash, cipherFactory func(key []byte) (cipher.Block, error), streamFactory func(cipher.Block, []byte) cipher.Stream) decrypter {
	// implementation note: the golang stdlib distinction between cipher.Stream and cipher.BlockMode is... odd, and could be readily wallpapered over with some really derpy wrappers.  haven't bothered yet.

	// also: feel slightly bad about passing in both the schema, and all its functors.  could make the Schema interface also return all these things.
	// really wish the hash/block/stream classes just returned their own sizes consistently, because that would make this a nonissue.

	// @implements decrypter
	return func(input io.Reader, output io.Writer, k Key) error {
		// read IV, use it to initialize cipher
		iv := make([]byte, sch.BlockSize())
		if _, err := io.ReadFull(input, iv); err != nil {
			return err
		}
		blockCipher, err := cipherFactory(k.cipherKey)
		if err != nil {
			return err
		}
		streamCipher := cipher.StreamWriter{
			S: streamFactory(blockCipher, iv),
			W: output,
		}

		// read length of ciphertext body
		var bodyLength int64
		if err := binary.Read(input, binary.BigEndian, &bodyLength); err != nil {
			return err
		}

		// read that much into hmac and the decipher stream
		hmacMsg := hmac.New(macFactory, k.hmacKey)
		mw := io.MultiWriter(streamCipher, hmacMsg)
		_, err = io.CopyN(mw, input, bodyLength)
		if err != nil {
			return err
		}

		// read and verify mac
		mac := make([]byte, sch.MACSize())
		if _, err := io.ReadFull(input, mac); err != nil {
			return err
		}
		if !hmac.Equal(mac, hmacMsg.Sum(nil)) {
			return err
			//return fmt.Errorf("unable to verify file")
		}

		return nil
	}
}
