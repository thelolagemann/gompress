package gompress

import (
	"testing"
)

func TestNewXz(t *testing.T) {
	l := NewXz()
	if l == nil {
		t.Error("NewXz returned nil pointer")
	}

}

func TestXz(t *testing.T) {
	c := NewXz()

	testCompressDecompress(t, c)
}

func BenchmarkXz(b *testing.B) {
	benchmarkCompressDecompress(b, NewXz())
}
