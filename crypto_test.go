package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"testing"
)

var (
	out [][]byte

	plaintext = mkRand(1024)
	keys      = []Key{
		Key{AES256_SHA256, mkRand(32), mkRand(32)},
		Key{AES256_Keccak256, mkRand(32), mkRand(32)},
		Key{Blowfish448_SHA256, mkRand(56), mkRand(32)},
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
