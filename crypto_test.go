package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"testing"
)

var (
	plaintextSize = 1024
	out           [][]byte

	plaintext = mkRand(plaintextSize)
	keys      = []Key{
		Key{AES256_SHA256, mkRand(AES256_SHA256.KeySize()), mkRand(AES256_SHA256.MACSize())},
		Key{AES256_Keccak256, mkRand(AES256_Keccak256.KeySize()), mkRand(AES256_Keccak256.MACSize())},
		Key{Blowfish448_SHA256, mkRand(Blowfish448_SHA256.KeySize()), mkRand(Blowfish448_SHA256.MACSize())},
	}
)

func mkRand(sz int) []byte {
	k := make([]byte, sz)
	io.ReadFull(rand.Reader, k)
	return k
}

func TestEncrypt(t *testing.T) {
	t.Logf("%20s: %.75s...\n", "plaintext", hex.EncodeToString(plaintext))
	for _, k := range keys {
		buf := new(bytes.Buffer)
		if err := Encrypt(bytes.NewBuffer(plaintext), buf, k); err != nil {
			t.Fatal(err)
		}
		t.Logf("%20s: %.75s...\n", k.Scheme, hex.EncodeToString(buf.Bytes()))
		out = append(out, buf.Bytes())
	}
	return
}

func TestDecrypt(t *testing.T) {
	t.Logf("%20s: %.75s...\n", "plaintext", hex.EncodeToString(plaintext))
	for i, k := range keys {
		x := new(bytes.Buffer)
		if err := Decrypt(bytes.NewBuffer(out[i]), x, k); err != nil {
			t.Fatal(err)
		}
		t.Logf("%20s: %.75s...\n", k.Scheme, hex.EncodeToString(x.Bytes()))
		if !bytes.Equal(plaintext, x.Bytes()) {
			t.Fail()
		}
	}
	return
}
