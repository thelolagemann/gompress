package gompress

import (
	"testing"
)

func TestNewLzma2(t *testing.T) {
	l := NewLzma2()
	if l == nil {
		t.Error("NewLzma2 returned nil pointer")
	}

}

func TestLzma2(t *testing.T) {
	c := NewLzma2()

	testCompressDecompress(t, c)
}

func BenchmarkLzma2(b *testing.B) {
	benchmarkCompressDecompress(b, NewLzma2())
}
