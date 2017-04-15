package receive

import (
	"cache"
	"constant"

	"github.com/spf13/viper"

	"fmt"
	"os"
	"time"
)

func NewSettings() (Settings, error) {
	var err error

	envS := os.Getenv("ENV")
	if envS == "" {
		envS = "development"
	}

	env, err := constant.ParseServerMode(envS)
	if err != nil {
		return Settings{}, err
	}

	v := viper.New()
	v.SetEnvPrefix("recive")
	v.BindEnv("deadline")
	v.SetDefault("deadline", 5*time.Second)

	conf := &Settings{}
	if err = v.Unmarshal(conf); err != nil {
		return *conf, fmt.Errorf("Invalid recive config: %v", err)
	}
	conf.Env = env

	conf.Cache, err = cache.NewSettings()

	return *conf, err
}

type Settings struct {
	Env             constant.ServerMode `mapstructure:"-"`
	DeadlineTimeout time.Duration       `mapstructure:"deadline"`

	Cache cache.Settings `mapstructure:"-"`
}
