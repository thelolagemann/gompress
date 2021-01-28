package gompress

import (
	"bytes"
	"testing"
)

func TestNewLz4(t *testing.T) {
	l := NewLz4()
	if l == nil {
		t.Error("NewLz4 returned nil")
	}

	if l.Level() != 5 {
		t.Errorf("expecting compression level of 5, got %v", l.level)
	}

	l = &Lz4{faultyCompressionLevel}
	l.SetLevel(11)
	if err := l.Compress(&bytes.Buffer{}, &bytes.Buffer{}); err == nil {
		t.Error("expecting an error compressing with invalid level, got none")
	}
}

func TestLz4(t *testing.T) {
	c := NewLz4()

	testCompressDecompress(t, c)
}

func BenchmarkLz4(b *testing.B) {
	benchmarkCompressDecompress(b, NewLz4())
}
