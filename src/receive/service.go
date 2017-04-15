package receive

import (
	"cache"

	"gopkg.in/gin-gonic/gin.v1"

	"context"
	"net/http"
)

func NewService(settings Settings) Receiver {
	return Receiver{settings}
}

type Receiver struct {
	settings Settings
}

func (s Receiver) Serve(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), s.settings.DeadlineTimeout)
	defer cancel()

	filename := c.Param("filename")

	res, err := s.handle(ctx, filename)
	if err != nil {
		if err == context.Canceled {
			return
		}

		var code int
		if err == context.DeadlineExceeded {
			code = http.StatusGatewayTimeout
			c.HTML(code, "504.html", gin.H{})
		} else {
			code = http.StatusNotFound
			c.HTML(code, "404.html", gin.H{})
		}
		c.AbortWithError(code, err)
		return
	}

	c.Header("Content-Type", res.ContentType)
	c.File(res.Filename)
}

func (s Receiver) handle(ctx context.Context, filename string) (*cache.File, error) {
	resCh := make(chan *cache.File)
	cacheService := cache.NewService(s.settings.Cache, ctx, resCh, filename)

	go cacheService.Serve()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-resCh:
		return result, result.Error()
	}
}
