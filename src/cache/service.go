package cache

import (
	"storage"

	"context"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func NewService(settings Settings, ctx context.Context, responseChan chan *File, filename string) Cache {
	return Cache{
		Filename: filename,

		ctx: ctx,

		responseChan: responseChan,

		settings: settings,
	}
}

type Cache struct {
	Filename string

	ctx context.Context

	responseChan chan *File

	settings Settings
}

// TODO заглушка
func (s Cache) Serve() {
	select {
	// Если контекст уже завершился, завершаем работу
	case <-s.ctx.Done():
		return
	default:
	}

	file := s.getFile()

	select {
	case <-s.ctx.Done():
	case s.responseChan <- file:
	}
}

func (s Cache) getFile() *File {
	db := s.settings.DB

	file := &File{
		Filename: filepath.Join(s.settings.CacheDir, url.PathEscape(s.Filename)),
	}

	if s.Filename == "/" {
		file.err = ErrDir
		return file
	}

	err := db.Where("filename = ?", file.Filename).Limit(1).Find(file).Error
	if err != nil {
		ff, err := s.getFromStorage()
		if err != nil {
			file.err = err
			return file
		}

		if err = file.Parse(ff); err != nil {
			file.err = err
			return file
		}

		if err = db.Create(file).Error; err != nil {
			file.err = err
			return file
		}

	} else {
		if file.Deleted {
			file.err = ErrDeleted
			return file
		}

		isDir, err := file.IsDir()
		if err != nil {
			file.err = err
			db.Model(file).UpdateColumn("deleted", true)
			return file
		}

		if isDir {
			file.err = ErrDir
			db.Model(file).UpdateColumn("deleted", true)
			return file
		}

		db.Model(file).UpdateColumn("updated_at", time.Now())
	}
	return file
}

func (s Cache) getFromStorage() (*os.File, error) {
	storageService := storage.NewService(s.settings.Storage, s.Filename)

	return storageService.Serve()
}
