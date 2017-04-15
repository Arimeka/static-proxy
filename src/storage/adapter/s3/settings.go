package s3

import (
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
	v.SetEnvPrefix("s3")
	v.BindEnv("id")
	v.BindEnv("bucket")
	v.BindEnv("secret")
	v.BindEnv("region")
	v.SetDefault("region", "eu-central-1")
	v.BindEnv("read_timeout")
	v.SetDefault("read_timeout", 5*time.Second)

	conf := &Settings{}
	if err = v.Unmarshal(conf); err != nil {
		return *conf, fmt.Errorf("Invalid s3 config: %v", err)
	}
	conf.Env = env

	return *conf, err
}

type Settings struct {
	Env    constant.ServerMode `mapstructure:"-"`
	ID     string
	Bucket string
	Secret string
	Region string

	ReadTimeout time.Duration `mapstructure:"read_timeout"`
}
