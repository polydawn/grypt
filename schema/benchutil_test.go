package schema

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"testing"
)

func encBench(b *testing.B, sch Schema, sz int) {
	k, err := sch.NewKey(rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	plaintext := bytes.NewReader(mkRand(sz))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := sch.Encrypt(plaintext, ioutil.Discard, k); err != nil {
			b.Fatal(err)
		}
		plaintext.Seek(0, 0)
	}
	return
}

func decBench(b *testing.B, sch Schema, sz int) {
	k, err := sch.NewKey(rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	buf := new(bytes.Buffer)
	buf.Grow(sz)
	if err := sch.Encrypt(bytes.NewReader(mkRand(sz)), buf, k); err != nil {
		b.Fatal(err)
	}
	r := bytes.NewReader(buf.Bytes())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := sch.Decrypt(r, ioutil.Discard, k); err != nil {
			b.Fatal(err)
		}
		r.Seek(0, 0)
	}
	return
}
