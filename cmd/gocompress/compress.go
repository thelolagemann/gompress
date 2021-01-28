package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/machinebox/progress"
	"github.com/schollz/progressbar/v3"
	"github.com/thelolagemann/gompress"
)

type compressCmd struct {
	Level  int    `short:"l" long:"level" description:"Compression level (if supported)"`
	Method string `short:"m" long:"method" description:"Compression method, use gompress list to view available methods" default:"gzip"`
	Output string `short:"o" long:"output" description:"Set the output destination"`
}

func (c *compressCmd) Execute(args []string) error {
	// get filecompressor
	fC, err := matchCompressor(c.Method)
	if err != nil {
		fatalLog(err)
	}

	if c.Level > 0 {
		l, leveller := fC.CompressorDecompressor.(gompress.Leveller)
		if !leveller {
			warnLog("%v doesn't support levels, ignoring", c.Method)
		} else {
			if err := l.SetLevel(c.Level); err != nil {
				errorLog(err)
			}
		}
	}

	for _, f := range args {
		var output string
		// TODO auto-generate output location
		if c.Output == "" {
			output = f + "." + fC.Extension
		}
		fmt.Printf("compressing %v => %v (%v)\n\n", f, output, c.Method)

		if err := compressFile(f, output, fC); err != nil {
			errorLog(err)
		}
	}

	return nil
}

// compressFile is a util function to compress a file
// located at input, compress with c, and write to output,
// displaying a progress bar as it progresses.
func compressFile(input string, output string, c gompress.CompressorDecompressor) error {

	// get input file
	info, err := os.Stat(input)
	if err != nil {
		return err
	}
	f, err := os.Open(input)
	if err != nil {
		return err
	}
	defer f.Close()

	// progress management
	pR := progress.NewReader(f)
	bar := progressbar.DefaultBytes(
		info.Size(),
		"compressing...",
	)

	// output management
	o, err := os.Create(output)
	if err != nil {
		return err
	}
	defer o.Close()
	pW := progress.NewWriter(o)
	go func() {
		ctx := context.Background()
		pTicker := progress.NewTicker(ctx, pR, info.Size(), 100*time.Millisecond)
		for p := range pTicker {
			// TODO calculate ratio
			bar.Set64(p.N())
		}
	}()

	if err := c.Compress(pR, pW); err != nil {
		return err
	}

	if err := bar.Finish(); err != nil {
		return err
	}

	return nil
}
