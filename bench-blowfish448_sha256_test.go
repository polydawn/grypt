package main

import "testing"

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
