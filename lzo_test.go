package gompress

import (
	"testing"
)

func TestNewLzo(t *testing.T) {
	l := NewLzo()
	if l == nil {
		t.Error("NewLzo returned nil pointer")
	}

}

func TestLzo(t *testing.T) {
	c := NewLzo()

	testCompressDecompress(t, c)
}

func BenchmarkLzo(b *testing.B) {
	benchmarkCompressDecompress(b, NewLzo())
}
