// +build !darwin

package gompress

import (
	"io"

	"github.com/cyberdelia/lzo"
)

// Lzo implements both the Compressor and Decompressor
// interfaces, facilitating the compression and decompression
// of data streams using the Lzo data format specification.
type Lzo struct{}

// NewLzo returns a new Lzo instance, ready to compress
// and decompress data streams.
func NewLzo() *Lzo {
	return &Lzo{}
}

// Compress implements the Compressor interface by wrapping
// w with Lzo.NewWriter, which is then used to copy data
// from r, compressing the data as it's copied.
func (l *Lzo) Compress(r io.Reader, w io.Writer) error {
	out := lzo.NewWriter(w)

	defer out.Close()

	_, err := io.Copy(out, r)
	return err
}

// Decompress implements the Decompressor interface by wrapping
// r with Lzo.NewReader, which is then used to copy data to w,
// decompressing the data as it's copied.
func (l *Lzo) Decompress(r io.Reader, w io.Writer) error {
	in, err := lzo.NewReader(r)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, in)

	return err
}
