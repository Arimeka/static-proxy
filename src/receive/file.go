package receive

import (
	"os"
	"errors"
)

type File struct {
	ContentType string
	Filename    string

	file        *os.File
}

func (f *File) Open() error {
	fi, err := os.Stat(f.Filename)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return errors.New("Path is directory")
	}

	ff, err := os.Open(f.Filename)
	if err != nil {
		return err
	}

	f.file = ff
	return nil
}

func (f *File) Close() error {
	return f.file.Close()
}

func (f *File) Delete() error {
	return os.Remove(f.Filename)
}
