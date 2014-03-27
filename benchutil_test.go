package main

import (
	"bytes"
	"crypto/rand"
	"testing"
)

var x int

func encBench(b *testing.B, s Scheme, sz int) {
	k, err := NewKey(rand.Reader, s)
	if err != nil {
		b.Fatal(err)
	}
	plaintext := bytes.NewReader(mkRand(sz))
	o := new(bytes.Buffer)
	o.Grow(sz)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Encrypt(plaintext, o, k); err != nil {
			b.Fatal(err)
		}
		x = o.Len()
		o.Reset()
	}
	return
}

func decBench(b *testing.B, s Scheme, sz int) {
	k, err := NewKey(rand.Reader, s)
	if err != nil {
		b.Fatal(err)
	}
	plaintext := bytes.NewBuffer(mkRand(sz))
	buf := new(bytes.Buffer)
	buf.Grow(sz)
	if err := Encrypt(plaintext, buf, k); err != nil {
		panic(err)
	}
	r := bytes.NewReader(buf.Bytes())
	o := new(bytes.Buffer)
	o.Grow(sz)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Decrypt(r, o, k); err != nil {
			b.Fatal(err)
		}
		r.Seek(0, 0)
		o.Reset()
	}
	return
}
