package gompress

import (
	"bytes"
	"testing"

	"github.com/klauspost/compress/zlib"
)

func TestNewZlib(t *testing.T) {
	z := NewZlib()
	if z == nil {
		t.Error("NewZlib returned nil")
	}

	if z.Level() != zlib.DefaultCompression {
		t.Errorf("expecting compression level of %v, got %v", zlib.DefaultCompression, z.level)
	}

	z = &Zlib{faultyCompressionLevel}
	z.SetLevel(zlib.BestCompression + 1)
	if err := z.Compress(&bytes.Buffer{}, &bytes.Buffer{}); err == nil {
		t.Error("expecting an error compressing data with an invalid level, got none")
	}

}

func TestZlib(t *testing.T) {
	c := NewZlib()

	testCompressDecompress(t, c)
}

func BenchmarkZlib(b *testing.B) {
	benchmarkCompressDecompress(b, NewZlib())
}
