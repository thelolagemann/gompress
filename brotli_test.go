package gompress

import (
	"testing"

	"github.com/andybalholm/brotli"
)

func TestNewBrotli(t *testing.T) {
	c := NewBrotli()
	if c == nil {
		t.Error("NewBrotli returned nil")
	}

	if c.level != brotli.DefaultCompression {
		t.Errorf("NewBrotli expecting compression level of %v, got %v", brotli.DefaultCompression, c.level)
	}

}

func TestBrotli(t *testing.T) {
	b := NewBrotli()
	testCompressDecompress(t, b)
}

func BenchmarkBrotli(b *testing.B) {
	benchmarkCompressDecompress(b, NewBrotli())
}
