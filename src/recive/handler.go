package recive

import (
	"github.com/Sirupsen/logrus"

	"context"
	"net/http"
	"time"
	"fmt"
)

func NewHandler(logger *logrus.Logger, dt time.Duration) Handler {
	return Handler{
		Log: logger,
		DeadlineTimeout: dt,
	}
}

type Handler struct {
	Log *logrus.Logger

	DeadlineTimeout time.Duration
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.DeadlineTimeout)
	defer cancel()

	res, err:= h.handle(ctx, w)
	if err != nil {
		http.Error(w, "Gateway Timeout", http.StatusGatewayTimeout)
		return
	}

	w.Write(res.Body)
}

func (h Handler) handle(ctx context.Context, w http.ResponseWriter) (*File, error) {
	resCh := make(chan []byte)
	h.Log.Info("Going to do hard work!")
	go hardWork(ctx, resCh)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
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
