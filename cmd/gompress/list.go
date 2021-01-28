package main

import (
	"fmt"

	"github.com/thelolagemann/gompress"
)

type listCmd struct{}

func (l *listCmd) Execute(args []string) error {
	fmt.Print("Methods\t\tSupported operations\tLevels\n\n")
	for _, c := range FileCompressors {
		fC, _ := c.(*gompress.FileCompressor)
		l, isLevelled := fC.CompressorDecompressor.(gompress.Leveller)
		if isLevelled {
			fmt.Printf("%v\t\tcompress,decompress\t%v - %v (default: %v)\n", fC.Method, l.MinLevel(), l.MaxLevel(), l.Level())
		} else {
			fmt.Printf("%v\t\tcompress,decompress\n", fC.Method)
		}
	}

	return nil
}
