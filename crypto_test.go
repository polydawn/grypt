package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Decrypt(r, ioutil.Discard, k); err != nil {
			b.Fatal(err)
		}
		r.Seek(0, 0)
	}
	return
}

func BenchmarkAES256_SHA256Encrypt1K(b *testing.B) {
	b.ReportAllocs()
	encBench(b, AES256_SHA256, 1024)
}
func BenchmarkAES256_SHA256Decrypt1K(b *testing.B) {
	b.ReportAllocs()
	decBench(b, AES256_SHA256, 1024)
}

func BenchmarkAES256_SHA256Encrypt4K(b *testing.B) {
	b.ReportAllocs()
	encBench(b, AES256_SHA256, 4*1024)
}
func BenchmarkAES256_SHA256Decrypt4K(b *testing.B) {
	b.ReportAllocs()
	decBench(b, AES256_SHA256, 4*1024)
}

func BenchmarkAES256_SHA256Encrypt1M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, AES256_SHA256, 1024*1024)
}
func BenchmarkAES256_SHA256Decrypt1M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, AES256_SHA256, 1024*1024)
}

func BenchmarkAES256_SHA256Encrypt2M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, AES256_SHA256, 2*(1024*1024))
}
func BenchmarkAES256_SHA256Decrypt2M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, AES256_SHA256, 2*(1024*1024))
}

func BenchmarkAES256_SHA256Encrypt4M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, AES256_SHA256, 4*(1024*1024))
}
func BenchmarkAES256_SHA256Decrypt4M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, AES256_SHA256, 4*(1024*1024))
}

func BenchmarkAES256_Keccak256Encrypt1K(b *testing.B) {
	b.ReportAllocs()
	encBench(b, AES256_Keccak256, 1024)
}
func BenchmarkAES256_Keccak256Decrypt1K(b *testing.B) {
	b.ReportAllocs()
	decBench(b, AES256_Keccak256, 1024)
}

func BenchmarkAES256_Keccak256Encrypt4K(b *testing.B) {
	b.ReportAllocs()
	encBench(b, AES256_Keccak256, 4*1024)
}
func BenchmarkAES256_Keccak256Decrypt4K(b *testing.B) {
	b.ReportAllocs()
	decBench(b, AES256_Keccak256, 4*1024)
}

func BenchmarkAES256_Keccak256Encrypt1M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, AES256_Keccak256, 1024*1024)
}
func BenchmarkAES256_Keccak256Decrypt1M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, AES256_Keccak256, 1024*1024)
}

func BenchmarkAES256_Keccak256Encrypt2M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, AES256_Keccak256, 2*(1024*1024))
}
func BenchmarkAES256_Keccak256Decrypt2M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, AES256_Keccak256, 2*(1024*1024))
}

func BenchmarkAES256_Keccak256Encrypt4M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, AES256_Keccak256, 4*(1024*1024))
}
func BenchmarkAES256_Keccak256Decrypt4M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, AES256_Keccak256, 4*(1024*1024))
}

func BenchmarkBlowfish448_SHA256Encrypt1K(b *testing.B) {
	b.ReportAllocs()
	encBench(b, Blowfish448_SHA256, 1024)
}
func BenchmarkBlowfish448_SHA256Decrypt1K(b *testing.B) {
	b.ReportAllocs()
	decBench(b, Blowfish448_SHA256, 1024)
}

func BenchmarkBlowfish448_SHA256Encrypt4K(b *testing.B) {
	b.ReportAllocs()
	encBench(b, Blowfish448_SHA256, 4*1024)
}
func BenchmarkBlowfish448_SHA256Decrypt4K(b *testing.B) {
	b.ReportAllocs()
	decBench(b, Blowfish448_SHA256, 4*1024)
}

func BenchmarkBlowfish448_SHA256Encrypt1M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, Blowfish448_SHA256, 1024*1024)
}
func BenchmarkBlowfish448_SHA256Decrypt1M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, Blowfish448_SHA256, 1024*1024)
}

func BenchmarkBlowfish448_SHA256Encrypt2M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, Blowfish448_SHA256, 2*(1024*1024))
}
func BenchmarkBlowfish448_SHA256Decrypt2M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, Blowfish448_SHA256, 2*(1024*1024))
}

func BenchmarkBlowfish448_SHA256Encrypt4M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, Blowfish448_SHA256, 4*(1024*1024))
}
func BenchmarkBlowfish448_SHA256Decrypt4M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, Blowfish448_SHA256, 4*(1024*1024))
}
