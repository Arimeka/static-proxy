package storage

import (
	"constant"
	"storage/adapter/s3"

	"fmt"
	"os"
)

func NewSettings() (Settings, error) {
	client, err := s3.NewClient()
	if err != nil {
		return Settings{}, fmt.Errorf("Invalid storage config: %v", err)
	}

	envS := os.Getenv("ENV")
	if envS == "" {
		envS = "development"
	}

	env, err := constant.ParseServerMode(envS)
	if err != nil {
		return Settings{}, err
	}

	return Settings{
		Env:    env,
		client: client,
	}, nil
}

type Settings struct {
	Env constant.ServerMode `mapstructure:"-"`

	client Client

	CacheDir string
}
