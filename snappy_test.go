package gompress

import (
	"testing"
)

func TestNewSnappy(t *testing.T) {
	s := NewSnappy()
	if s == nil {
		t.Error("NewSnappy returned nil")
	}
}

func TestSnappy(t *testing.T) {
	c := NewSnappy()

	testCompressDecompress(t, c)
}

func BenchmarkSnappy(b *testing.B) {
	benchmarkCompressDecompress(b, NewSnappy())
}
