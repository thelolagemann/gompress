package gompress

import (
	"io"

	"github.com/andybalholm/brotli"
)

// Brotli implements both the Compressor and Decompressor
// interfaces, facilitating the compression and decompression
// of data streams using the brotli data format specification.
type Brotli struct {
	*CompressionLevel
}

// NewBrotli returns a new Brotli instance, ready to compress
// and decompress data streams.
func NewBrotli() *Brotli {
	return &Brotli{&CompressionLevel{
		level:    brotli.DefaultCompression,
		minLevel: brotli.BestSpeed + 1,
		maxLevel: 11,
	}}
}

// Compress implements the Compressor interface by wrapping
// w with brotli.NewWriterLevel, which is then used to copy
// data from r, compressing the data as it's copied.
func (br *Brotli) Compress(r io.Reader, w io.Writer) error {
	out := brotli.NewWriterLevel(w, br.level)
	defer out.Close()
	_, err := io.Copy(out, r)
	return err
}

// Decompress implements the Decompressor interface by wrapping
// r with brotli.NewReader, which is then used to copy data from r
// compressing the data as it's copied.
func (br *Brotli) Decompress(r io.Reader, w io.Writer) error {
	in := brotli.NewReader(r)

	_, err := io.Copy(w, in)
	return err
}
