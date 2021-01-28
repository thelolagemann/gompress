/*
Package gompress provides both a CLI interface, and a golang library
for working with various compression algorithms. Currently
supported are

	brotli
	bzip2
	gzip
	lz4
	lzma2
	snappy
	zlib
	zstd


*/
package gompress

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	errNoCompressorDecompressor            = errors.New("FileCompressor initialized without CompressorDecompressor")
	fs                          fileSystem = &osFS{}
)

type fileSystem interface {
	Create(filename string) (*os.File, error)
	Open(filename string) (*os.File, error)
	Stat(filename string) (os.FileInfo, error)
}

type osFS struct {
}

func (o osFS) Create(filename string) (*os.File, error)  { return os.Create(filename) }
func (o osFS) Open(filename string) (*os.File, error)    { return os.Open(filename) }
func (o osFS) Stat(filename string) (os.FileInfo, error) { return os.Stat(filename) }

// Compressor is a type that is capable of compressing
// a data stream provided by an io.Reader, writing the
// output to an io.Writer.
type Compressor interface {
	// Compress compresses data from r, and writes to w.
	Compress(r io.Reader, w io.Writer) error
}

// Decompressor is a type that is capable of decompressing
// a data stream provided by an io.Reader, writing the
// output to an io.Writer.
type Decompressor interface {
	// Decompress decompresses data from r, and writes to w.
	Decompress(r io.Reader, w io.Writer) error
}

type CompressorDecompressor interface {
	Compressor
	Decompressor
}

// Leveller defines the contracts for compression methods
// that make use of levels should implement in order to allow
// for the control of the internal compression algorithms
// compression level.
type Leveller interface {
	Level() int
	MinLevel() int
	MaxLevel() int
	SetLevel(level int) error
}

// FileCompressor serves as a wrapper for Compressor/Decompressors,
// providing util functions for working with files.
type FileCompressor struct {
	CompressorDecompressor
	// Default file extension
	Extension string
	// The name of the compression method
	Method string
}

// CompressionLevel is a helper struct to allow for easy implementation
// of the Leveller interface.
type CompressionLevel struct {
	level    int
	minLevel int
	maxLevel int
}

// Level returns the compression level
func (c *CompressionLevel) Level() int {
	return c.level
}

// MinLevel returns the minimum supported compression level.
func (c *CompressionLevel) MinLevel() int {
	return c.minLevel
}

// MaxLevel returns the maximum supported compression level.
func (c *CompressionLevel) MaxLevel() int {
	return c.maxLevel
}

// SetLevel sets the compression level to l if it's within
// the bounds of c.MinLevel() and c.MaxLevel(), returning
// an error otherwise.
func (c *CompressionLevel) SetLevel(l int) error {
	if inBound(l, c.minLevel, c.maxLevel) {
		c.level = l
		return nil
	}

	return fmt.Errorf("level %v out of range, supported levels: %v - %v", l, c.minLevel, c.maxLevel)
}

// NewFileCompressor returns a new FileCompressor
// instance, ready to compress and decompress files.
func NewFileCompressor(c CompressorDecompressor, method, extension string) *FileCompressor {
	fC := &FileCompressor{
		CompressorDecompressor: c,
		Extension:              extension,
		Method:                 method,
	}
	return fC
}

// CompressFile reads the file from input and writes the
// compressed data to dest, using the Compress method defined
// by the Compressor interface.
func (f *FileCompressor) CompressFile(input string, dest string) error {
	if f.CompressorDecompressor == nil {
		return errNoCompressorDecompressor
	}
	if isExists(dest) {
		return os.ErrExist
	}
	r, err := fs.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := fs.Create(dest)
	if err != nil {
		return err
	}

	defer w.Close()

	return f.Compress(r, w)
}

// DecompressFile reads the compressed file from input
// and writes the decompressed data to dest, using the
// Decompress method defined by the Decompressor interface.
func (f *FileCompressor) DecompressFile(input string, dest string) error {
	if f.CompressorDecompressor == nil {
		return errNoCompressorDecompressor
	}
	if isExists(dest) {
		return os.ErrExist
	}

	r, err := fs.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := fs.Create(dest)
	if err != nil {
		return err
	}
	defer w.Close()

	return f.Decompress(r, w)
}

func isExists(file string) bool {
	_, err := fs.Stat(file)
	return !os.IsNotExist(err)
}
