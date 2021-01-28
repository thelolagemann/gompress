package gompress

import (
	"io"

	"github.com/golang/snappy"
)

// Snappy implements both the Compressor and Decompressor
// interfaces, facilitating the compression and decompression
// of data streams using the snappy data compression algorithm.
type Snappy struct{}

// NewSnappy returns a new Snappy instance, ready to compress
// and decompress data streams.
func NewSnappy() *Snappy {
	return &Snappy{}
}

// Compress implements the Compressor interface by wrapping
// w with snappy.NewBufferedWriter, which is then used to copy
// data from r, compressing the data as it's copied.
func (s *Snappy) Compress(r io.Reader, w io.Writer) error {
	out := snappy.NewBufferedWriter(w)
	defer out.Close()
	_, err := io.Copy(out, r)
	return err
}

// Decompress implements the Decompressor interface by wrapping
// r with snappy.NewReader, which is then used to copy data to w,
// decompressing the data as it's copied.
func (s *Snappy) Decompress(r io.Reader, w io.Writer) error {
	in := snappy.NewReader(r)
	_, err := io.Copy(w, in)
	return err
}
