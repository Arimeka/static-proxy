package settings

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"strings"
)

type Settings struct {
	Address    string
	Port       string
	Workers    int
	Scheme     string
	S3Config   *S3Config
	ValidSizes *ValidSizes
}

type S3Config struct {
	Hosts map[string]map[string]string
}

type ValidSizes struct {
	Sizes map[string][]string
}

var (
	Config      *Settings    = &Settings{}
	AppSettings *viper.Viper = viper.New()
	S3Settings  *viper.Viper = viper.New()
)

func SetupOptions(options string) {
	if strings.Contains(options, "config") {
		flag.StringP("config", "c", "", "path to config file")
	}
	if strings.Contains(options, "env") {
		flag.StringP("env", "e", "development", "environment")
	}
	if strings.Contains(options, "address") {
		flag.StringP("bind", "b", "", "bind address")
		flag.StringP("port", "p", "", "bind port")
	}
}

func Setup() {
	flag.Parse()

	AppSettings.BindPFlag("port", flag.Lookup("port"))
	AppSettings.BindPFlag("address", flag.Lookup("bind"))
	AppSettings.BindPFlag("config", flag.Lookup("config"))
	AppSettings.BindPFlag("env", flag.Lookup("env"))
	S3Settings.BindPFlag("env", flag.Lookup("env"))

	AppSettings.SetDefault("address", "0.0.0.0")
	AppSettings.SetDefault("port", "5000")
	AppSettings.SetDefault("env", "development")
	AppSettings.SetDefault("config", "./config.yml")
	S3Settings.SetDefault("env", "development")

	configPath := filepath.Dir(AppSettings.GetString("config"))
	configName := strings.Replace(filepath.Base(AppSettings.GetString("config")), filepath.Ext(AppSettings.GetString("config")), "", -1)

	AppSettings.SetConfigName(configName)
	AppSettings.AddConfigPath(configPath)

	err := AppSettings.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	S3Settings.SetConfigName("s3")
	S3Settings.AddConfigPath("./")

	err = S3Settings.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func Build() (err error) {

	err = AppSettings.UnmarshalKey(AppSettings.GetString("env"), Config)
	if err != nil {
		return
	}

	Config.Address = AppSettings.GetString("address")
	Config.Port = AppSettings.GetString("port")

	s3Config := &S3Config{}
	err = S3Settings.UnmarshalKey(S3Settings.GetString("env"), s3Config)
	if err != nil {
		return
	}
	Config.S3Config = s3Config

	sizes := viper.New()
	sizes.SetConfigName("sizes")
	sizes.AddConfigPath("./")

	err = sizes.ReadInConfig()
	if err != nil {
		return
	}

	mapSizes := &ValidSizes{}
	err = sizes.UnmarshalKey(AppSettings.GetString("env"), mapSizes)
	if err != nil {
		return
	}
	Config.ValidSizes = mapSizes

	return
}
