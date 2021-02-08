package gompress

import (
	"io"

	"github.com/pierrec/lz4/v4"
)

// Lz4 implements both the Compressor and Decompressor
// interfaces, facilitating the compression and decompression
// of data streams using the lz4 data compression algorithm.
type Lz4 struct {
	*CompressionLevel
}

// NewLz4 returns a new Lz4 instance, ready to compress
// and decompress data streams.
func NewLz4() *Lz4 {
	return &Lz4{&CompressionLevel{
		level:    5,
		minLevel: 1,
		maxLevel: 9,
	}}
}

// Compress implements the Compressor interface by wrapping
// w with lz4.NewWriter, which is then used to copy data
// from r, compressing the data as it's copied.
func (l *Lz4) Compress(r io.Reader, w io.Writer) error {
	out := lz4.NewWriter(w)
	options := []lz4.Option{
		lz4.CompressionLevelOption(lz4.CompressionLevel(1 << (8 + l.level))),
	}
	defer out.Close()
	if err := out.Apply(options...); err != nil {
		return err
	}

	_, err := io.Copy(out, r)
	return err
}

// Decompress implements the Decompressor interface by wrapping
// r with lz4.NewReader, which is then used to copy data from r
// compressing the data as it's copied.
func (l *Lz4) Decompress(r io.Reader, w io.Writer) error {
	in := lz4.NewReader(r)

	_, err := io.Copy(w, in)
	return err
}
