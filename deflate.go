package gompress

import (
	"compress/flate"
	"io"
)

// Deflate implements both the Compressor and Decompressor
// interfaces, facilitating the compression and decompression
// of data streams using the Deflate data format specification.
type Deflate struct {
	*CompressionLevel
}

// NewDeflate returns a new Deflate instance, ready to compress
// and decompress data streams.
func NewDeflate() *Deflate {
	return &Deflate{&CompressionLevel{
		level:    flate.DefaultCompression,
		minLevel: flate.BestSpeed,
		maxLevel: flate.BestCompression,
	}}
}

// Compress implements the Compressor interface by wrapping
// w with Deflate.NewWriter, which is then used to copy data
// from r, compressing the data as it's copied.
func (b *Deflate) Compress(r io.Reader, w io.Writer) error {
	out, err := flate.NewWriter(w, b.level)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, r)
	return err
}

// Decompress implements the Decompressor interface by wrapping
// r with Deflate.NewReader, which is then used to copy data to w,
// decompressing the data as it's copied.
func (b *Deflate) Decompress(r io.Reader, w io.Writer) error {
	// Deflate.NewReader never returns error, so no need to check
	in := flate.NewReader(r)

	defer in.Close()
	_, err := io.Copy(w, in)
	return err
}
