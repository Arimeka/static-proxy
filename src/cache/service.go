package cache

import (
	"context"
	"net/url"
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

	db := s.settings.DB

	file := &File{
		Filename: filepath.Join(s.settings.CacheDir, url.PathEscape(s.Filename)),
	}

	err := db.Where("filename = ?", file.Filename).Limit(1).Find(file).Error
	if err != nil {
		file.err = err
	} else if file.Deleted {
		file.err = ErrDeleted
	} else {
		isDir, err := file.IsDir()
		if err != nil {
			file.err = err
			db.Model(file).UpdateColumn("deleted", true)
		} else if isDir {
			file.err = ErrDir
			db.Model(file).UpdateColumn("deleted", true)
		} else {
			db.Model(file).UpdateColumn("updated_at", time.Now())
		}
	}

	select {
	case <-s.ctx.Done():
	case s.responseChan <- file:
	}
}
