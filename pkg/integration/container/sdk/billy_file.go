package sdk

import (
	"github.com/go-git/go-billy/v5"
	"os"
)

type billyFile struct {
	f    billy.File
	stat os.FileInfo
}

func (f billyFile) Close() error {
	if f.f != nil {
		return f.f.Close()
	}

	return nil
}

func (f billyFile) Read(buf []byte) (int, error) {
	if f.f != nil {
		return f.f.Read(buf)
	}

	return 0, ErrNotSupported
}

func (f billyFile) Write(buf []byte) (int, error) {
	if f.f != nil {
		return f.f.Write(buf)
	}

	return 0, ErrNotSupported
}

func (f billyFile) Seek(offset int64, whence int) (int64, error) {
	if f.f != nil {
		return f.f.Seek(offset, whence)
	}

	return 0, ErrNotSupported
}

// Readdir does...something. I have no idea what. This implementation actually
// just always returns an ErrNotSupported error because I couldn't figure out
// what behavior for this function should actually be.
func (f billyFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, ErrNotSupported
}

func (f billyFile) Stat() (os.FileInfo, error) {
	return f.stat, nil
}
