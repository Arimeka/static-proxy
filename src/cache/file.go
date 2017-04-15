package cache

import (
	"errors"
	"mime"
	"os"
	"time"
)

var (
	ErrDir     = errors.New("Path is directory")
	ErrDeleted = errors.New("File marked as deleted")
)

type File struct {
	ID uint `gorm:"primary_key"`

	ContentType string
	Size        int64
	Filename    string

	Deleted bool

	CreatedAt time.Time
	UpdatedAt time.Time

	err error `gorm:"-"`

	file *os.File `gorm:"-"`
}

func (f *File) Parse(ff *os.File) error {
	fi, err := os.Stat(ff.Name())
	if err != nil {
		return err
	}

	f.file = ff
	f.Size = fi.Size()
	f.ContentType = mime.TypeByExtension(f.Filename)

	return nil
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
