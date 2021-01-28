package gompress

import (
	"errors"
	"io"

	"github.com/ulikunitz/xz/lzma"
)

// Lzma2 implements both the Compressor and Decompressor
// interfaces, facilitating the compression and decompression
// of data streams using the lzma2 data format specification.
type Lzma2 struct{}

// NewLzma2 returns a new Lzma2 instance, ready to compress
// and decompress data streams.
func NewLzma2() *Lzma2 {
	return &Lzma2{}
}

// Compress implements the Compressor interface by wrapping
// w with lzma.NewWriter2, which is then used to copy data
// from r, compressing the data as it's copied.
func (l *Lzma2) Compress(r io.Reader, w io.Writer) error {
	// err only when initializing from lzma.WriterConfig
	out, _ := lzma.NewWriter2(w)
	defer out.Close()

	_, err := io.Copy(out, r)
	return err
}

// Decompress implements the Decompressor interface by wrapping
// r with lzma.NewReader2, which is then used to copy data to w,
// decompressing the data as it's copied.
func (l *Lzma2) Decompress(r io.Reader, w io.Writer) error {
	// err only returned when initializing from lzma.ReaderConfig2
	in, _ := lzma.NewReader2(r)

	n, err := io.Copy(w, in)
	// error doesn't happen when copying, check n
	if n == 0 {
		return errors.New("invalid lzma2 data")
	}
	return err
}
