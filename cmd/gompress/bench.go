package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/thelolagemann/gompress"
)

type benchCmd struct {
	All bool `short:"a" long:"all" description:"Benchmark all methods"`
}

func (b *benchCmd) Execute(args []string) error {

	if b.All {
		fmt.Println("benchmarking all compression methods, this may take a while...")
		for _, f := range args {
			for _, c := range FileCompressors {
				if err := benchmarkFile(f, c); err != nil {
					fatalLog(err)
				}
			}

		}
	}
	return nil
}

type benchmarkResult struct {
	timeTaken      time.Time
	compressedSize int64
}

// benchmarkFile is a util function to benchmark a file
// located at input, running through each of the compression
// methods, reporting the time taken and resulting ratio.
func benchmarkFile(input string, c gompress.CompressorDecompressor) error {
	// get input
	info, err := os.Stat(input)
	if err != nil {
		return err
	}
	f, err := os.Open(input)
	if err != nil {
		return err
	}
	defer f.Close()
	fC, _ := c.(*gompress.FileCompressor)
	l, hasLevel := fC.CompressorDecompressor.(gompress.Leveller)
	if hasLevel {
		for n := l.MinLevel(); n <= l.MaxLevel(); n++ {
			f.Seek(0, 0)
			start := time.Now()
			if err := l.SetLevel(n); err != nil {
				return err
			}
			var b bytes.Buffer
			if err := c.Compress(f, &b); err != nil {
				return err
			}
			ratio := float64(b.Len()) / float64(info.Size()) * 100
			fmt.Printf("%v (level: %v)\t%v (%.2f%%)\t%v\n", fC.Method, n, len(b.Bytes()), ratio, time.Since(start))
		}
	} else {
		var b bytes.Buffer
		start := time.Now()

		if err := fC.Compress(f, &b); err != nil {
			return err
		}
		ratio := float64(b.Len()) / float64(info.Size()) * 100
		fmt.Printf("%v\t%v (%.2f%%)\t%v\n", fC.Method, len(b.Bytes()), ratio, time.Since(start))
	}

	return nil
}
