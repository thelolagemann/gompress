// +build !darwin

package main

import "github.com/thelolagemann/gompress"

func init() {
	FileCompressors = append(FileCompressors, gompress.NewFileCompressor(gompress.NewLzo(), "lzo", "lzo"))
}
