package cache

import (
	"context"
	"mime"
	"net/url"
	"path/filepath"
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

	file := &File{
		Filename: filepath.Join("./cache", url.PathEscape(s.Filename)),
	}
	file.ContentType = mime.TypeByExtension(file.Filename)

	isDir, err := file.IsDir()
	if err != nil {
		file.err = err
	} else if isDir {
		file.err = ErrDir
	}

	select {
	case <-s.ctx.Done():
	case s.responseChan <- file:
	}
}
