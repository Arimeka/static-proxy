package transform

import (
	"constant"

	"github.com/spf13/viper"

	"fmt"
	"os"
	"strings"
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
	v.SetEnvPrefix("transfromt")
	v.BindEnv("allowed_sizes")
	v.SetDefault("allowed_sizes", "250x250,250x211,640x430,847x410,320x380,600x,x564,560x")
	v.BindEnv("max_pixels")
	v.SetDefault("max_pixels", 8000)

	conf := &Settings{}
	if err := v.Unmarshal(conf); err != nil {
		return *conf, fmt.Errorf("Invalid tranfsorm config: %v", err)
	}
	conf.Env = env
	conf.AllowedSizes = strings.Split(v.GetString("allowed_sizes"), ",")

	return *conf, nil
}

type Settings struct {
	Env constant.ServerMode `mapstructure:"-"`

	AllowedSizes []string `mapstructure:"-"`
	MaxPixels    uint64   `mapstructure:"max_pixels"`

	CacheDir string `mapstructure:"-"`
}
