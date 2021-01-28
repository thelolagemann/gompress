package gompress

import (
	"compress/flate"
	"io"

	"github.com/klauspost/compress/gzip"
)

// Gzip implements both the Compressor and Decompressor
// interfaces, facilitating the compression and decompression
// of data streams using the gzip data format specification.
type Gzip struct {
	*CompressionLevel
}

// NewGzip returns a new Gzip instance, ready to compress
// and decompress data streams.
func NewGzip() *Gzip {
	return &Gzip{&CompressionLevel{
		level:    flate.DefaultCompression,
		minLevel: flate.BestSpeed,
		maxLevel: flate.BestCompression,
	}}
}

// Compress implements the Compressor interface by wrapping
// w with gzip.NewWriterLevel, which is then used to copy data
// from r, compressing the data as it's copied.
func (g *Gzip) Compress(r io.Reader, w io.Writer) error {
	out, err := gzip.NewWriterLevel(w, g.level)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, r)
	return err
}

// Decompress implements the Decompressor interface by wrapping
// r with gzip.NewReader, which is then used to copy data to w,
// decompressing the data as it's copied.
func (g *Gzip) Decompress(r io.Reader, w io.Writer) error {
	in, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer in.Close()
	_, err = io.Copy(w, in)
	return err
}
