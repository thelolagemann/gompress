package gompress

import (
	"io"

	"github.com/ulikunitz/xz"
)

// Xz implements both the Compressor and Decompressor
// interfaces, facilitating the compression and decompression
// of data streams using the Xz data format specification.
type Xz struct{}

// NewXz returns a new Xz instance, ready to compress
// and decompress data streams.
func NewXz() *Xz {
	return &Xz{}
}

// Compress implements the Compressor interface by wrapping
// w with xz.NewWriter, which is then used to copy data
// from r, compressing the data as it's copied.
func (x *Xz) Compress(r io.Reader, w io.Writer) error {
	out, _ := xz.NewWriter(w)
	defer out.Close()

	_, err := io.Copy(out, r)
	return err
}

// Decompress implements the Decompressor interface by wrapping
// r with xz.NewReader, which is then used to copy data to w,
// decompressing the data as it's copied.
func (x *Xz) Decompress(r io.Reader, w io.Writer) error {
	in, err := xz.NewReader(r)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, in)

	return err
}
