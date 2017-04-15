package storage

import (
	"context"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

func NewService(settings Settings, filename string) Storage {
	return Storage{
		Filename: filename,
		settings: settings,
	}
}

type Storage struct {
	Filename string

	settings Settings
}

func (s Storage) Serve() (*os.File, error) {
	reader, cancel, err := s.settings.client.Get(s.Filename)
	if err != nil {
		return nil, err
	}
	defer cancel()

	file, err := os.Create(filepath.Join(s.settings.CacheDir, url.PathEscape(s.Filename)))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		os.Remove(file.Name())
		if err == context.Canceled {
			err = errors.New("Storage cancel request")
		}
		return nil, err
	}
	if err = reader.Close(); err != nil {
		os.Remove(file.Name())
		if err == context.Canceled {
			err = errors.New("Storage cancel request")
		}
		return nil, err
	}

	return file, nil
}
