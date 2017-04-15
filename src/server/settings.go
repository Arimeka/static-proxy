package main

import (
	"constant"
	"receive"

	"github.com/spf13/viper"

	"fmt"
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

	envS := os.Getenv("ENV")
	if envS == "" {
		envS = "development"
	}

	env, err := constant.ParseServerMode(envS)
	if err != nil {
		return Settings{}, err
	}

	v := viper.New()
	v.SetEnvPrefix("server")
	v.BindEnv("addr")
	v.SetDefault("addr", ":5000")
	v.BindEnv("readtimeout")
	v.SetDefault("readtimeout", 5*time.Second)
	v.BindEnv("writetimeout")
	v.SetDefault("writetimeout", 15*time.Second)

	conf := &Settings{}
	if err = v.Unmarshal(conf); err != nil {
		return *conf, fmt.Errorf("Invalid server config: %v", err)
	}
	conf.Env = env
	conf.WrkDir = wrkDir

	conf.Receiever, err = receive.NewSettings()

	return *conf, err
}

type Settings struct {
	Env        constant.ServerMode `mapstructure:"-"`
	WrkDir     string              `mapstructure:"-"`
	ServerAddr string              `mapstructure:"addr"`

	ReadTimeout  time.Duration `mapstructure:"readtimeout"`
	WriteTimeout time.Duration `mapstructure:"writetimeout"`

	Receiever receive.Settings `mapstructure:"-"`
}
