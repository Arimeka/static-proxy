package receive

import (
	"gopkg.in/gin-gonic/gin.v1"

	"context"
	"net/http"
	"time"
	"errors"
	"mime"
	"path/filepath"
)

func NewService(dt time.Duration) Receiver {
	return Receiver{
		deadlineTimeout: dt,
	}
}

type Receiver struct {
	deadlineTimeout time.Duration
}

func (s Receiver) Serve(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), s.deadlineTimeout)
	defer cancel()

	cn, ok := c.Writer.(http.CloseNotifier)
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("Bad Request"))
		return
	}

	filename := c.Param("filename")

	res, err:= s.handle(ctx, cn, filename)
	if err != nil {
		var code int
		if err == context.DeadlineExceeded {
			code = http.StatusGatewayTimeout
			c.HTML(code, "504.html",gin.H{})
		} else {
			code = http.StatusNotFound
			c.HTML(code, "404.html",gin.H{})
		}
		c.AbortWithError(code, err)
		return
	}

	c.File(res.Filename)
}

func (s Receiver) handle(ctx context.Context, cn http.CloseNotifier, filename string) (*File, error) {
	resCh := make(chan *File)
	errCh := make(chan error)
	go hardWork(ctx, resCh, errCh, filename)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-cn.CloseNotify():
		return nil, errors.New("client closed connection")
	case err := <-errCh:
		return nil, err
	case result := <-resCh:
		return result, nil
	}
}

// TODO заглушка
func hardWork(ctx context.Context, responseChan chan *File, errorChan chan error , filename string) {
	select {
	// Если контекст уже завершился, завершаем работу
	case <-ctx.Done():
		return
	default:
	}

	file := &File{
		Filename: filepath.Join("./cache", filename),
	}
	file.ContentType = mime.TypeByExtension(file.Filename)

	err := file.Open()
	if err!= nil {
		errorChan <- err
		return
	}
	file.Close()

	select {
	case <-ctx.Done():
	case responseChan <- file:
	}
}
