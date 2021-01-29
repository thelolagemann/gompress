package gompress

import (
	"io"

	"github.com/ulikunitz/xz/lzma"
)

// Lzma implements both the Compressor and Decompressor
// interfaces, facilitating the compression and decompression
// of data streams using the lzma data format specification.
type Lzma struct{}

// NewLzma returns a new Lzma instance, ready to compress
// and decompress data streams.
func NewLzma() *Lzma {
	return &Lzma{}
}

// Compress implements the Compressor interface by wrapping
// w with lzma.NewWriter, which is then used to copy data
// from r, compressing the data as it's copied.
func (l *Lzma) Compress(r io.Reader, w io.Writer) error {
	out, _ := lzma.NewWriter(w)

	defer out.Close()

	_, err := io.Copy(out, r)
	return err
}

// Decompress implements the Decompressor interface by wrapping
// r with lzma.NewReader, which is then used to copy data to w,
// decompressing the data as it's copied.
func (l *Lzma) Decompress(r io.Reader, w io.Writer) error {
	in, err := lzma.NewReader(r)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, in)

	return err
}
