package gompress

import "testing"

func TestNewLzma(t *testing.T) {
	l := NewLzma()
	if l == nil {
		t.Error("NewLzma returned nil pointer")
	}
}

func TestLzma(t *testing.T) {
	l := NewLzma()

	testCompressDecompress(t, l)
}

func BenchmarkLzma(b *testing.B) {
	benchmarkCompressDecompress(b, NewLzma())
}
