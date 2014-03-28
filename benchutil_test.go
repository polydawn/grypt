package main

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"testing"
)

func encBench(b *testing.B, s Scheme, sz int) {
	k, err := NewKey(rand.Reader, s)
	if err != nil {
		b.Fatal(err)
	}
	plaintext := bytes.NewReader(mkRand(sz))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Encrypt(plaintext, ioutil.Discard, k); err != nil {
			b.Fatal(err)
		}
		plaintext.Seek(0, 0)
	}
	return
}

func decBench(b *testing.B, s Scheme, sz int) {
	k, err := NewKey(rand.Reader, s)
	if err != nil {
		b.Fatal(err)
	}
	buf := new(bytes.Buffer)
	buf.Grow(sz)
	if err := Encrypt(bytes.NewReader(mkRand(sz)), buf, k); err != nil {
		b.Fatal(err)
	}
	r := bytes.NewReader(buf.Bytes())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Decrypt(r, ioutil.Discard, k); err != nil {
			b.Fatal(err)
		}
		r.Seek(0, 0)
	}
	return
}
