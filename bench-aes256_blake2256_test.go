package main

import "testing"

func BenchmarkAES256_BLAKE2256Encrypt1K(b *testing.B) {
	b.ReportAllocs()
	encBench(b, AES256_BLAKE2256, 1024)
}
func BenchmarkAES256_BLAKE2256Decrypt1K(b *testing.B) {
	b.ReportAllocs()
	decBench(b, AES256_BLAKE2256, 1024)
}

func BenchmarkAES256_BLAKE2256Encrypt4K(b *testing.B) {
	b.ReportAllocs()
	encBench(b, AES256_BLAKE2256, 4*1024)
}
func BenchmarkAES256_BLAKE2256Decrypt4K(b *testing.B) {
	b.ReportAllocs()
	decBench(b, AES256_BLAKE2256, 4*1024)
}

func BenchmarkAES256_BLAKE2256Encrypt1M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, AES256_BLAKE2256, 1024*1024)
}
func BenchmarkAES256_BLAKE2256Decrypt1M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, AES256_BLAKE2256, 1024*1024)
}

func BenchmarkAES256_BLAKE2256Encrypt2M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, AES256_BLAKE2256, 2*(1024*1024))
}
func BenchmarkAES256_BLAKE2256Decrypt2M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, AES256_BLAKE2256, 2*(1024*1024))
}

func BenchmarkAES256_BLAKE2256Encrypt4M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, AES256_BLAKE2256, 4*(1024*1024))
}
func BenchmarkAES256_BLAKE2256Decrypt4M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, AES256_BLAKE2256, 4*(1024*1024))
}
