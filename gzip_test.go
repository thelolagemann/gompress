package gompress

import (
	"bytes"
	"compress/gzip"
	"testing"
)

func TestNewGzip(t *testing.T) {
	c := NewGzip()
	if c == nil {
		t.Error("NewGzip returned nil pointer")
	}
	if c.Level() != gzip.DefaultCompression {
		t.Errorf("expecting compression level of %v, got %v", gzip.DefaultCompression, c.level)
	}

	c = &Gzip{faultyCompressionLevel}
	c.SetLevel(gzip.BestCompression + 10)
	if err := c.Compress(&bytes.Buffer{}, &bytes.Buffer{}); err == nil {
		t.Error("expecting an error compressing with an invalid level, got none")
	}

}

func TestGzip(t *testing.T) {
	c := NewGzip()
	testCompressDecompress(t, c)
}

func BenchmarkGzip(b *testing.B) {
	benchmarkCompressDecompress(b, NewGzip())
}
