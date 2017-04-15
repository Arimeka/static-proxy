package receive

import (
	"cache"
	"time"
)

type Settings struct {
	DeadlineTimeout time.Duration

	Cache cache.Settings
}
