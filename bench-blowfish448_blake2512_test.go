package main

import "testing"

func BenchmarkBlowfish448_BLAKE2512Encrypt1K(b *testing.B) {
	b.ReportAllocs()
	encBench(b, Blowfish448_BLAKE2512, 1024)
}
func BenchmarkBlowfish448_BLAKE2512Decrypt1K(b *testing.B) {
	b.ReportAllocs()
	decBench(b, Blowfish448_BLAKE2512, 1024)
}

func BenchmarkBlowfish448_BLAKE2512Encrypt4K(b *testing.B) {
	b.ReportAllocs()
	encBench(b, Blowfish448_BLAKE2512, 4*1024)
}
func BenchmarkBlowfish448_BLAKE2512Decrypt4K(b *testing.B) {
	b.ReportAllocs()
	decBench(b, Blowfish448_BLAKE2512, 4*1024)
}

func BenchmarkBlowfish448_BLAKE2512Encrypt1M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, Blowfish448_BLAKE2512, 1024*1024)
}
func BenchmarkBlowfish448_BLAKE2512Decrypt1M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, Blowfish448_BLAKE2512, 1024*1024)
}

func BenchmarkBlowfish448_BLAKE2512Encrypt2M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, Blowfish448_BLAKE2512, 2*(1024*1024))
}
func BenchmarkBlowfish448_BLAKE2512Decrypt2M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, Blowfish448_BLAKE2512, 2*(1024*1024))
}

func BenchmarkBlowfish448_BLAKE2512Encrypt4M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, Blowfish448_BLAKE2512, 4*(1024*1024))
}
func BenchmarkBlowfish448_BLAKE2512Decrypt4M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, Blowfish448_BLAKE2512, 4*(1024*1024))
}
