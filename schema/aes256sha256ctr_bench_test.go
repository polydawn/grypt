package schema

import "testing"

func BenchmarkAES256_SHA256Encrypt1K(b *testing.B) {
	b.ReportAllocs()
	encBench(b, Aes256sha256ctr{}, 1024)
}
func BenchmarkAES256_SHA256Decrypt1K(b *testing.B) {
	b.ReportAllocs()
	decBench(b, Aes256sha256ctr{}, 1024)
}

func BenchmarkAES256_SHA256Encrypt4K(b *testing.B) {
	b.ReportAllocs()
	encBench(b, Aes256sha256ctr{}, 4*1024)
}
func BenchmarkAES256_SHA256Decrypt4K(b *testing.B) {
	b.ReportAllocs()
	decBench(b, Aes256sha256ctr{}, 4*1024)
}

func BenchmarkAES256_SHA256Encrypt1M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, Aes256sha256ctr{}, 1024*1024)
}
func BenchmarkAES256_SHA256Decrypt1M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, Aes256sha256ctr{}, 1024*1024)
}

func BenchmarkAES256_SHA256Encrypt2M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, Aes256sha256ctr{}, 2*(1024*1024))
}
func BenchmarkAES256_SHA256Decrypt2M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, Aes256sha256ctr{}, 2*(1024*1024))
}

func BenchmarkAES256_SHA256Encrypt4M(b *testing.B) {
	b.ReportAllocs()
	encBench(b, Aes256sha256ctr{}, 4*(1024*1024))
}
func BenchmarkAES256_SHA256Decrypt4M(b *testing.B) {
	b.ReportAllocs()
	decBench(b, Aes256sha256ctr{}, 4*(1024*1024))
}
