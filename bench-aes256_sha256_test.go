package main

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"testing"

	"polydawn.net/grypt/schema"
)

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

func BenchmarkNewAES256_SHA256Encrypt1K(b *testing.B) {
	schem := schema.Aes256sha256ctr{}
	b.ReportAllocs()
	k, err := schem.NewKey(rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	plaintext := bytes.NewReader(mkRand(1024))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := schem.Encrypt(plaintext, ioutil.Discard, k); err != nil {
			b.Fatal(err)
		}
		plaintext.Seek(0, 0)
	}
	return
}
