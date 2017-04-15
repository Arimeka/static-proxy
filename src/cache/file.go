package cache

import (
	"errors"
	"os"
	"time"
)

var ErrDir = errors.New("Path is directory")

type File struct {
	ID uint `gorm:"primary_key"`

	ContentType string
	Size        int64
	Filename    string

	CreatedAt time.Time
	UpdatedAt time.Time

	err error `gorm:"-"`

	file *os.File `gorm:"-"`
}

func (f *File) Open() error {
	isDir, err := f.IsDir()
	if err != nil {
		return err
	} else if isDir {
		return ErrDir
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

func (f File) IsDir() (bool, error) {
	fi, err := os.Stat(f.Filename)
	if err != nil {
		return false, err
	}

	return fi.IsDir(), nil
}

func (f File) Error() error {
	return f.err
}
