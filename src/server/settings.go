package main

import (
	"receive"

	"cache"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func NewSettings() (Settings, error) {
	var wrkDir string
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return Settings{}, err
	}

	dirs := strings.Split(dir, "/")
	index := len(dirs) - 1
	if dirs[index] == "bin" {
		dirs = append(dirs[:index], dirs[index+1:]...)
		wrkDir = strings.Join(dirs, "/")
	} else {
		wrkDir = dir
	}

	return Settings{
		WrkDir:       wrkDir,
		ServerAddr:   ":5000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		Receiever: receive.Settings{
			DeadlineTimeout: 5 * time.Second,

			Cache: cache.Settings{
				StorageLimit: 15 << (10 * 2), // 15 MB
			},
		},
	}, nil
}

type Settings struct {
	WrkDir     string
	ServerAddr string

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	Receiever receive.Settings
}
