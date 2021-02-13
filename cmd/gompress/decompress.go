package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thelolagemann/gompress"
)

type decompressCmd struct {
	Method string `short:"m" long:"method" description:"Override the inferred compressor"`
}

func (d *decompressCmd) Execute(args []string) error {
	var fC *gompress.FileCompressor
	var err error

	// determine file compression
	if d.Method == "" {
		// get decompressor by extension
		fC, err = matchCompressor(filepath.Ext(args[0])[1:])
	} else {
		// get decompressor by supplied method
		fC, err = matchCompressor(d.Method)

	}
	if err != nil {
		return err
	}

	for _, f := range args {
		output := strings.Replace(f, "."+fC.Extension, "", 1)

		fmt.Printf("decompressing %v => %v (%v)\n\n", f, output, fC.Method)

		if err := decompressFile(f, output, fC); err != nil {
			errorLog(err)
		}
	}

	return nil
}

func decompressFile(input string, output string, d gompress.CompressorDecompressor) error {
	// get input file
	f, err := os.Open(input)
	if err != nil {
		return err
	}
	defer f.Close()

	// output
	o, err := os.Create(output)
	if err != nil {
		return err
	}
	defer o.Close()
	if err := d.Decompress(f, o); err != nil {
		return err
	}

	return nil
}
