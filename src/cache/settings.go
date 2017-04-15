package cache

import (
	"constant"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	"fmt"
	"os"
	"storage"
	"time"
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
	v.SetDefault("limit", 15*constant.MB)
	v.BindEnv("dir")
	v.SetDefault("dir", "cache")
	v.BindEnv("clear_duration")
	v.SetDefault("clear_duration", 15*time.Minute)

	conf := &Settings{}
	if err := v.Unmarshal(conf); err != nil {
		return *conf, fmt.Errorf("Invalid cache config: %v", err)
	}
	conf.Env = env

	conf.Storage, err = storage.NewSettings()
	if err != nil {
		return *conf, err
	}
	conf.Storage.CacheDir = conf.CacheDir

	db, err := NewDBConn()
	if err != nil {
		return *conf, fmt.Errorf("Failed connected to cache DB: %v", err)
	}
	conf.DB = db

	stats, err := NewStats(db)
	if err != nil {
		return *conf, fmt.Errorf("Failed calculate cache stats: %v", err)
	}
	stats.LimitSize = conf.StorageLimit
	stats.CleanDuration = conf.ClearDuration
	conf.Stats = stats

	timer := time.NewTimer(stats.CleanDuration)
	go stats.CacheWatcher(timer)

	return *conf, nil
}

type Settings struct {
	Env          constant.ServerMode `mapstructure:"-"`
	StorageLimit constant.ByteSize   `mapstructure:"limit"`
	CacheDir     string              `mapstructure:"dir"`

	ClearDuration time.Duration `mapstructure:"clear_duration"`

	Storage storage.Settings `mapstructure:"-"`
	DB      *gorm.DB         `mapstructure:"-"`
	Stats   *Stats           `mapstructure:"-"`
}
