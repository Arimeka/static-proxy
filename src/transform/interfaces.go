package transform

import "os"

type Worker interface {
	Process() (*os.File, error)
}
