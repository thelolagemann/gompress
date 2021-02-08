package gompress

import (
	"io"

	"github.com/dsnet/compress/bzip2"
)

// Bzip2 implements both the Compressor and Decompressor
// interfaces, facilitating the compression and decompression
// of data streams using the bzip2 data format specification.
type Bzip2 struct {
	*CompressionLevel
}

// NewBzip2 returns a new Bzip2 instance, ready to compress
// and decompress data streams.
func NewBzip2() *Bzip2 {
	return &Bzip2{&CompressionLevel{
		level:    bzip2.DefaultCompression,
		minLevel: bzip2.BestSpeed,
		maxLevel: bzip2.BestCompression,
	}}
}

// Compress implements the Compressor interface by wrapping
// w with bzip2.NewWriter, which is then used to copy data
// from r, compressing the data as it's copied.
func (b *Bzip2) Compress(r io.Reader, w io.Writer) error {
	out, err := bzip2.NewWriter(w, &bzip2.WriterConfig{
		Level: b.level,
	})
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, r)
	return err
}

// Decompress implements the Decompressor interface by wrapping
// r with bzip2.NewReader, which is then used to copy data to w,
// decompressing the data as it's copied.
func (b *Bzip2) Decompress(r io.Reader, w io.Writer) error {
	// bzip2.NewReader never returns error, so no need to check
	in, _ := bzip2.NewReader(r, nil)

	defer in.Close()
	_, err := io.Copy(w, in)
	return err
}
