package gompress

import (
	"compress/flate"
	"io"

	"github.com/klauspost/compress/zlib"
)

// Zlib implements both the Compressor and Decompressor
// interfaces, facilitating the compression and decompression
// of data streams using the zlib data format specification.
type Zlib struct {
	*CompressionLevel
}

// NewZlib returns a new Zlib instance, ready to compress
// and decompress data streams.
func NewZlib() *Zlib {
	return &Zlib{&CompressionLevel{
		level:    flate.DefaultCompression,
		minLevel: flate.BestSpeed,
		maxLevel: flate.BestCompression,
	}}
}

// Compress implements the Compressor interface by wrapping
// w with zlib.NewWriter, which is then used to copy data
// from r, compressing the data as it's copied.
func (z *Zlib) Compress(r io.Reader, w io.Writer) error {
	out, err := zlib.NewWriterLevel(w, z.level)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, r)
	return err
}

// Decompress implements the Decompressor interface by wrapping
// r with zlib.NewReader, which is then used to copy data to w,
// decompressing the data as it's copied.
func (z *Zlib) Decompress(r io.Reader, w io.Writer) error {
	in, err := zlib.NewReader(r)
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.Copy(w, in)
	return err
}
