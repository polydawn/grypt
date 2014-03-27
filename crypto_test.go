package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"testing"
)

var (
	plaintext []byte
	out       [][]byte

	keys = []Key{
		Key{AES256_SHA256, []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")},
		Key{AES256_Keccak256, []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")},
	}
)

func init() {
	plaintext = make([]byte, 1024)
	io.ReadFull(rand.Reader, plaintext)
}

func TestEncrypt(t *testing.T) {
	t.Logf("plaintext:\t%.75s...\n", hex.EncodeToString(plaintext))
	for _, k := range keys {
		buf := new(bytes.Buffer)
		if err := Encrypt(bytes.NewBuffer(plaintext), buf, k); err != nil {
			t.Fatal(err)
		}
		t.Logf("ciphertext:\t%.75s...\n", hex.EncodeToString(buf.Bytes()))
		out = append(out, buf.Bytes())
	}
	return
}

func TestDecrypt(t *testing.T) {
	t.Logf("known plaintext:\t%.75s...\n", hex.EncodeToString(plaintext))
	for i, k := range keys {
		x := new(bytes.Buffer)
		if err := Decrypt(bytes.NewBuffer(out[i]), x, k); err != nil {
			t.Fatal(err)
		}
		t.Logf("output plaintext:\t%.75s...\n", hex.EncodeToString(x.Bytes()))
		if !bytes.Equal(plaintext, x.Bytes()) {
			t.Fail()
		}
	}
	return
}

func BenchmarkEncrypt(b *testing.B) {
	return
}

func BenchmarkDecrypt(b *testing.B) {
	return
}
