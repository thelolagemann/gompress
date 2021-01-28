package gompress

import (
	"bytes"
	"testing"

	"github.com/dsnet/compress/bzip2"
)

func TestNewBzip2(t *testing.T) {
	c := NewBzip2()
	if c == nil {
		t.Error("NewBzip2Compressor returned nil")
	}

	// check compression level
	if c.Level() != bzip2.DefaultCompression {
		t.Errorf("NewBzip2Compressor expected compression level %v, got %v", bzip2.DefaultCompression, c.level)
	}
}

func TestBzip2(t *testing.T) {
	c := NewBzip2()

	testCompressDecompress(t, c)

	c = &Bzip2{faultyCompressionLevel}
	c.SetLevel(bzip2.BestCompression + 1)
	if err := c.Compress(&bytes.Buffer{}, &bytes.Buffer{}); err == nil {
		t.Error("expecting an error compressing with an invalid level, got none")
	}
}

func BenchmarkBzip2(b *testing.B) {
	benchmarkCompressDecompress(b, NewBzip2())
}
