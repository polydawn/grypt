package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"testing"
	"testing/iotest"
)

var (
	plaintextSize = 1024
	out           [][]byte

	plaintext = mkRand(plaintextSize)
	keys      = []Key{
		Key{AES256_SHA256, mkRand(AES256_SHA256.KeySize()), mkRand(AES256_SHA256.MACSize())},
		Key{AES256_Keccak256, mkRand(AES256_Keccak256.KeySize()), mkRand(AES256_Keccak256.MACSize())},
		Key{Blowfish448_SHA256, mkRand(Blowfish448_SHA256.KeySize()), mkRand(Blowfish448_SHA256.MACSize())},
		Key{AES256_BLAKE2256, mkRand(AES256_BLAKE2256.KeySize()), mkRand(AES256_BLAKE2256.MACSize())},
		Key{Blowfish448_BLAKE2512, mkRand(Blowfish448_BLAKE2512.KeySize()), mkRand(Blowfish448_BLAKE2512.MACSize())},
	}
)

func mkRand(sz int) []byte {
	k := make([]byte, sz)
	io.ReadFull(rand.Reader, k)
	return k
}

func TestEncrypt(t *testing.T) {
	t.Logf("%25s: %.75s...\n", "plaintext", hex.EncodeToString(plaintext))
	for _, k := range keys {
		buf := new(bytes.Buffer)
		if err := Encrypt(bytes.NewReader(plaintext), buf, k); err != nil {
			t.Fatal(err)
		}
		t.Logf("%25s: %.75s...\n", k.Scheme, hex.EncodeToString(buf.Bytes()))
		out = append(out, buf.Bytes())
	}
	return
}

func TestDecrypt(t *testing.T) {
	t.Logf("%25s: %.75s...\n", "plaintext", hex.EncodeToString(plaintext))
	for i, k := range keys {
		x := new(bytes.Buffer)
		if err := Decrypt(bytes.NewReader(out[i]), x, k); err != nil {
			t.Fatal(err)
		}
		t.Logf("%25s: %.75s...\n", k.Scheme, hex.EncodeToString(x.Bytes()))
		if !bytes.Equal(plaintext, x.Bytes()) {
			t.Fail()
		}
	}
	return
}

func TestMACFailure(t *testing.T) {
	for _, k := range keys {
		var err error
		buf := new(bytes.Buffer)
		if err := Encrypt(bytes.NewReader(plaintext), iotest.TruncateWriter(buf, int64(plaintextSize-2)), k); err != nil {
			t.Fatal(err)
		}
		if err = Decrypt(buf, ioutil.Discard, k); err == nil {
			t.Logf("This should have errored! %s", k.Scheme)
			t.Fail()
		} else {
			t.Logf("%25s: %v\n", k.Scheme, err)
		}
	}
	return
}
