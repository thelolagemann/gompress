package gompress

import (
	"testing"
)

func TestNewZstd(t *testing.T) {
	c := NewZstd()
	if c == nil {
		t.Error("NewZstd returned nil")
	}
}

func TestZstd(t *testing.T) {
	c := NewZstd()

	testCompressDecompress(t, c)
}

func BenchmarkZstd(b *testing.B) {
	benchmarkCompressDecompress(b, NewZstd())
}
