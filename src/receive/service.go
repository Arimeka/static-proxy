package receive

import (
	"gopkg.in/gin-gonic/gin.v1"

	"context"
	"net/http"
	"time"
	"errors"
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


	res, err:= s.handle(ctx, cn)
	if err != nil {
		c.AbortWithError(http.StatusGatewayTimeout, errors.New("Gateway Timeout"))
		return
	}

	c.String(http.StatusOK,"%s",res.Body)
}

func (s Receiver) handle(ctx context.Context, cn http.CloseNotifier) (*File, error) {
	resCh := make(chan []byte)
	go hardWork(ctx, resCh)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-cn.CloseNotify():
		return nil, errors.New("client closed connection")
	case result := <-resCh:
		return &File{result}, nil
	}
}

// TODO заглушка
func hardWork(ctx context.Context, responseChan chan []byte) error {
	select {
	// Если контекст уже завершился, завершаем работу
	case <-ctx.Done():
		return nil
	default:
	}

	//time.Sleep(10*time.Second)
	hardJobResult := []byte("Welcome!\n")

	select {
	case <-ctx.Done():
	case responseChan <- hardJobResult:
	}
	return nil
}
