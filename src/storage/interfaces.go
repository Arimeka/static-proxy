package storage

import (
	"context"
	"io"
)

type Client interface {
	Get(filename string) (io.ReadCloser, context.CancelFunc, error)
}
