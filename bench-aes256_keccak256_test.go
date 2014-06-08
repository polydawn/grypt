package main

import "testing"

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
