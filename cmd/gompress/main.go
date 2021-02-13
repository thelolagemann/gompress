/*
gompress provides a CLI interface for
the gompress package.

Usage: gompress (bench|compress|decompress|list|help) [arguments...]

bench
	Benchmarks all available compression methods,
	reporting the time taken, and the resulting
	output size.

	- all
		Benchmark all algorithms (default: true)
	-i|input
		Input file to benchmark (default: "./gompress")
compress
	Compresses an individual file.

	-method
		Compression method (default: "gzip")
	-i|input
		Input file to compress
	-o|output
		Output destination (optional)
decompress
	Decompresses an individual file.

	- i|input
		Input file to decompress
	- o|output
		Output destination (optional)
list
	List all available compression methods,
	and their capabilities.
*/
package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/thelolagemann/gompress"
)

var (
	FileCompressors []gompress.CompressorDecompressor = []gompress.CompressorDecompressor{
		gompress.NewFileCompressor(gompress.NewBrotli(), "brotli", "br"),
		gompress.NewFileCompressor(gompress.NewBzip2(), "bzip2", "bz2"),
		gompress.NewFileCompressor(gompress.NewDeflate(), "deflate", "z"),
		gompress.NewFileCompressor(gompress.NewGzip(), "gzip", "gz"),
		gompress.NewFileCompressor(gompress.NewLzma(), "lzma", "lzma"),
		gompress.NewFileCompressor(gompress.NewLzma2(), "lzma2", "lz2"),
		gompress.NewFileCompressor(gompress.NewLz4(), "lz4", "lz4"),
		gompress.NewFileCompressor(gompress.NewSnappy(), "snappy", "snz"),
		gompress.NewFileCompressor(gompress.NewXz(), "xz", "xz"),
		gompress.NewFileCompressor(gompress.NewZlib(), "zlib", "zlib"),
		gompress.NewFileCompressor(gompress.NewZstd(), "zstd", "zst"),
	}
	infoLog  = func(message string, args ...interface{}) { fmt.Print("[ info  ] "+message, args) }
	warnLog  = func(message string, args ...interface{}) { fmt.Print("[ warn  ] "+message, args) }
	errorLog = func(err error, args ...interface{}) { fmt.Printf("[ error ] %v\n", err.Error()) }
	fatalLog = func(err error, args ...interface{}) {
		fmt.Printf("[ fatal ] %v\n", err)
		os.Exit(1)
	}

	bCmd benchCmd
	cCmd compressCmd
	dCmd decompressCmd
	lCmd listCmd

	opts Options
)

type Options struct {
	Verbose bool `short:"v" long:"verbose" description:"Verbose output"`
}

// TODO benchmarks

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.AddCommand("benchmark",
		"Benchmark compression methods",
		"The benchmark command will report time taken, and compression ratio for the selected compression methods. Use -a|all to benchmark all methods, this will take a while!",
		&bCmd)
	parser.AddCommand("compress",
		"Compress an individual file",
		`The compress command will compress an individual file, if no output is specified, 
		then gompress will automatically infer the output based on the compression method used.
		Compression level may be specified with the -l|level switch if it's supported.`,
		&cCmd)
	parser.AddCommand("decompress",
		"Decompress an individual file",
		`The decompress command will decompress an individual file. You may override the automatically
		inferred compression method by specifying -m|method`,
		&dCmd)
	parser.AddCommand("list",
		"List available compression methods",
		"The list command will list all the available compression methods and if applicable, their supported levels",
		&lCmd)
	parser.Parse()
}

// matchCompressor is a util function that iterates
// through FileCompressors trying to match the Method
// or Extension field of the FileCompressor. If no matches
// are found, an error will be returned.
func matchCompressor(match string) (*gompress.FileCompressor, error) {
	for _, c := range FileCompressors {
		fC, _ := c.(*gompress.FileCompressor)
		if fC.Method == match || fC.Extension == match {
			return fC, nil
		}
	}

	return nil, fmt.Errorf("no file compressor found matching %v", match)
}

func humanizeBytes2(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
