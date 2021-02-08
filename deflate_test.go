package gompress

import (
	"bytes"
	"compress/flate"
	"testing"
)

func TestNewDeflate(t *testing.T) {
	c := NewDeflate()
	if c == nil {
		t.Error("NewDeflateCompressor returned nil")
	}

	// check compression level
	if c.Level() != flate.DefaultCompression {
		t.Errorf("NewDeflateCompressor expected compression level %v, got %v", flate.DefaultCompression, c.level)
	}
}

func TestDeflate(t *testing.T) {
	c := NewDeflate()

	testCompressDecompress(t, c)

	c = &Deflate{faultyCompressionLevel}
	c.SetLevel(flate.BestCompression + 1)
	if err := c.Compress(&bytes.Buffer{}, &bytes.Buffer{}); err == nil {
		t.Error("expecting an error compressing with an invalid level, got none")
	}
}

func BenchmarkDeflate(b *testing.B) {
	benchmarkCompressDecompress(b, NewDeflate())
}
