package gompress

import (
	"errors"
	"io"

	"github.com/klauspost/compress/zstd"
)

// Zstd implements both the Compressor and Decompressor
// interfaces, facilitating the compression and decompression
// of data streams using the zstd data format specification.
type Zstd struct{}

// NewZstd returns a new Zstd instance, ready to compress
// and decompress data streams.
func NewZstd() *Zstd {
	return &Zstd{}
}

// Compress implements the Compressor interface by wrapping
// w with zstd.NewWriter, which is then used to copy data
// from r, compressing the data as it's copied.
func (z *Zstd) Compress(r io.Reader, w io.Writer) error {
	// err only returned when applying options
	out, _ := zstd.NewWriter(w)
	defer out.Close()

	_, err := io.Copy(out, r)
	return err
}

// Decompress implements the Decompressor interface by wrapping
// r with zstd.NewReader, which is then used to copy data to w,
// decompressing the data as it's copied.
func (z *Zstd) Decompress(r io.Reader, w io.Writer) error {
	// err only returned when applying options
	in, _ := zstd.NewReader(r)
	defer in.Close()

	n, err := io.Copy(w, in)
	if n == 0 {
		return errors.New("invalid zstd data")
	}
	return err
}
