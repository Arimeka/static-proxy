package cache

import (
	"constant"

	"github.com/spf13/viper"

	"fmt"
	"os"
)

func NewSettings() (Settings, error) {
	envS := os.Getenv("ENV")
	if envS == "" {
		envS = "development"
	}

	env, err := constant.ParseServerMode(envS)
	if err != nil {
		return Settings{}, err
	}

	v := viper.New()
	v.SetEnvPrefix("cache")
	v.BindEnv("limit")
	v.SetDefault("limit", 15<<(10*2)) // 15 MB

	conf := &Settings{}
	if err := v.Unmarshal(conf); err != nil {
		return *conf, fmt.Errorf("Invalid cache config: %v", err)
	}
	conf.Env = env

	return *conf, nil
}

type Settings struct {
	Env          constant.ServerMode `mapstructure:"-"`
	StorageLimit uint64              `mapstructure:"limit"`
}
