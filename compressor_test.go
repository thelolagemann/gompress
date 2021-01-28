package gompress

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
	"testing"
)

const (
	testFolder = "./testdata"
	testFile   = "./testdata/enwik9"
)

var (
	m1                = 1024 * 1024
	defaultBenchmarks = []struct {
		name string
		size int
	}{
		{"4M", m1 * 4},
		{"16M", m1 * 16},
		{"64M", m1 * 64},
		{"256M", m1 * 256},
	}
	defaultTest = make([]byte, 1024*1024)

	// used to test invalid compression levels
	faultyCompressionLevel = &CompressionLevel{
		level:    5,
		minLevel: 1,
		maxLevel: 100,
	}

	errMock = errors.New("mock error")
)

type mockCompressor struct{}

func (m *mockCompressor) Compress(r io.Reader, w io.Writer) error {
	return nil
}

func (m *mockCompressor) Decompress(r io.Reader, w io.Writer) error {
	return nil
}

type errMockCompressor struct{}

func (e *errMockCompressor) Compress(r io.Reader, w io.Writer) error {
	return errMock
}

func (e *errMockCompressor) Decompress(r io.Reader, w io.Writer) error {
	return errMock
}

type mockFs struct {
	osFS
	errCreate error
	errOpen   error
	errStat   error
}

type mockedFileInfo struct {
	os.FileInfo
}

func (m mockFs) Create(filename string) (*os.File, error) {
	if m.errCreate != nil {
		return nil, m.errCreate
	}

	return &os.File{}, nil
}

func (m mockFs) Open(filename string) (*os.File, error) {
	if m.errOpen != nil {
		return nil, m.errOpen
	}

	return &os.File{}, nil
}

func (m mockFs) Stat(filename string) (os.FileInfo, error) {
	if m.errStat != nil {
		return nil, m.errStat
	}

	return &mockedFileInfo{}, nil
}

func (m *mockFs) reset() {
	m.errCreate = nil
	m.errOpen = nil
	m.errStat = os.ErrNotExist // default for isExists
}

func TestNewFileCompressor(t *testing.T) {
	fC := NewFileCompressor(&mockCompressor{}, "mock", "mock")

	if fC == nil {
		t.Error("NewFileCompressor returned nil pointer")
	}
}

func TestFileCompressor(t *testing.T) {
	fC := NewFileCompressor(&mockCompressor{}, "mock", "mock")
	oldFs := fs
	mfs := &mockFs{}
	fs = mfs
	defer func() {
		fs = oldFs
	}()

	mfs.errStat = os.ErrNotExist
	// test successful CompressFile/DecompressFile, set errStat to handle isExists
	t.Run("CompressFile", func(t *testing.T) {
		if err := fC.CompressFile(testFile, testFile+".mock"); err != nil {
			t.Error(err)
		}
	})
	t.Run("DecompressFile", func(t *testing.T) {
		if err := fC.DecompressFile(testFile+".mock", testFile+".tmp"); err != nil {
			t.Error(err)
		}
	})

	// test existing destination output, set errStat to os.ErrExist
	mfs.errStat = os.ErrExist
	t.Run("ExistingFileDest", func(t *testing.T) {
		if err := fC.CompressFile(testFile, testFile); err != os.ErrExist {
			t.Errorf("expecting os.ErrExist, got %v", err)
		}
		if err := fC.DecompressFile(testFile, testFile); err != os.ErrExist {
			t.Errorf("expecting os.ErrExist, got %v", err)
		}
	})

	// test non-existing file input, set errOpen to os.ErrNotExist
	mfs.reset()
	mfs.errOpen = os.ErrNotExist
	t.Run("NonExistingFileInput", func(t *testing.T) {
		if err := fC.CompressFile(testFolder+"/invalid", testFolder+"/invalid.mock"); err != os.ErrNotExist {
			t.Errorf("expecting os.ErrNotExist, got %v", err)
		}
		if err := fC.DecompressFile(testFolder+"/invalid", testFolder+"/invalid"); err != os.ErrNotExist {
			t.Errorf("expecting os.ErrNotExist, got %v", err)
		}
	})

	// test error creating output file
	mfs.reset()
	mfs.errCreate = errMock
	t.Run("ErrorCreatingFileOutput", func(t *testing.T) {
		if err := fC.CompressFile(testFile, testFile+".mock"); err != errMock {
			t.Errorf("expecting errMock, got %v", err)
		}
		if err := fC.DecompressFile(testFile+".mock", testFile); err != errMock {
			t.Errorf("expecting errMock, got %v", err)
		}
	})

	// test FileCompressor initialized without CompressorDecompressor
	t.Run("NoCompressorDecompressor", func(t *testing.T) {
		fC := &FileCompressor{}
		if err := fC.CompressFile(testFile, testFile+".mock"); err != errNoCompressorDecompressor {
			t.Errorf("expecting errNoCompressorDecompressor, got %v", err)
		}
		if err := fC.DecompressFile(testFile+".mock", testFile); err != errNoCompressorDecompressor {
			t.Errorf("expecting errNoCompressorDecompressor, got %v", err)
		}
	})

	// test CompressorDecompressor returning error
	mfs.reset()
	t.Run("ErrCompressorDecompressor", func(t *testing.T) {
		fC := NewFileCompressor(&errMockCompressor{}, "err-mock", "err-mock")
		if err := fC.CompressFile(testFile, testFile+".mock"); err == nil {
			t.Errorf("expecting errMock compressing file, got %v", err)
		}

		if err := fC.DecompressFile(testFile+".mock", testFile+".tmp"); err == nil {
			t.Errorf("expecting errMock decompressing file, got %v", err)
		}
	})
}

func testCompressDecompress(t *testing.T, cD CompressorDecompressor) {
	var in bytes.Buffer
	var out bytes.Buffer

	var buffer = make([]byte, 1024*512)
	readTestData(buffer)
	in.Write(buffer)

	t.Run("Compress", func(t *testing.T) {
		if err := cD.Compress(&in, &out); err != nil {
			t.Error(err)
		}
	})

	t.Run("Decompress", func(t *testing.T) {
		var compare bytes.Buffer
		if err := cD.Decompress(&out, &compare); err != nil {
			t.Error(err)
		}

		if bytes.Compare(compare.Bytes(), buffer) != 0 {
			t.Error("decompression output was not the same as compression input")
		}
	})

	t.Run("DecompressInvalidData", func(t *testing.T) {
		var in bytes.Buffer
		var out bytes.Buffer

		var buffer = make([]byte, 1024*512)
		readTestData(buffer)
		in.Write(buffer)
		if err := cD.Decompress(&in, &out); err == nil {
			// some compressors don't return error for decompress, and simply copy data
			if len(in.Bytes()) == 0 || len(out.Bytes()) == 0 && bytes.Compare(in.Bytes(), out.Bytes()) != 0 {
				t.Error("expecting an error, or for data to copy, neither case happened")
			}
		}
	})

	l, leveller := cD.(Leveller)
	if leveller {
		t.Run("Level", func(t *testing.T) {
			if err := l.SetLevel(l.MinLevel() - 1); err == nil {
				t.Error("expecting an error whilst setting invalid compression level, got none")
			}
			if err := l.SetLevel(l.MaxLevel() + 1); err == nil {
				t.Error("expecting an error whilst setting invalid compression level, got none")
			}
		})

	}
}

func benchmarkCompressDecompress(b *testing.B, cD CompressorDecompressor) {
	for _, bench := range defaultBenchmarks {
		var w bytes.Buffer
		var r bytes.Buffer
		b.Run(bench.name, func(b *testing.B) {
			var buffer = make([]byte, bench.size)
			readTestData(buffer)
			r.Write(buffer)
			if !testing.Short() {
				l, isLevelled := cD.(Leveller)
				if isLevelled {
					for n := l.MinLevel(); n <= l.MaxLevel(); n++ {
						b.Run("CompressLevel-"+strconv.Itoa(n), func(b *testing.B) {
							for i := 0; i < b.N; i++ {
								if err := cD.Compress(&r, &w); err != nil {
									b.Error(err)
								}
							}
						})
					}
				}
			} else {
				var buffer = make([]byte, 1024*1024)
				r.Write(buffer)
				for i := 0; i < b.N; i++ {
					if err := cD.Compress(&r, &w); err != nil {
						b.Error(err)
					}
				}
			}
		})
	}
}
