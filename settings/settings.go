package settings

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
)

var DEFAULTS = map[string]interface{}{
	"env":     "dev",
	"address": "0.0.0.0",
	"port":    "5000",
	"scheme":  "http",
	"workers": 5,
}

type Settings struct {
	Address string
	Port    string
	NumCPU  int    `yaml:"numcpu"`
	Workers int    `yaml:"workers"`
	Scheme  string `yaml:"scheme"`
}

type S3 struct {
	Host map[interface{}]interface{} `yaml:"hosts"`
}

var (
	Env     string
	Address string
	Port    string
	Path    string
	S3Path  string
	Workers int
)

func SetupOptions(options string) {
	if strings.Contains(options, "config") {
		flag.StringVar(&Path, "config", "", "path to config file")
		flag.StringVar(&Path, "c", "", "path to config file (short)")
		flag.StringVar(&S3Path, "s3", "", "path to s3 config file")
		flag.StringVar(&S3Path, "s", "", "path to s3 config file (short)")
	}
	if strings.Contains(options, "env") {
		flag.StringVar(&Env, "environment", "development", "environment")
		flag.StringVar(&Env, "e", "development", "environment (short)")
		flag.IntVar(&Workers, "workers", 0, "number of workers")
		flag.IntVar(&Workers, "w", 0, "number op workers (short)")
	}
	if strings.Contains(options, "address") {
		flag.StringVar(&Address, "bind", "", "bind address")
		flag.StringVar(&Address, "b", "", "bind address (short)")
		flag.StringVar(&Port, "port", "", "bind port")
		flag.StringVar(&Port, "p", "", "bind port (short)")
	}
}

func Setup() {
	flag.Parse()

	if Env == "" {
		Env = DEFAULTS["env"].(string)
	}
}

func SetDefaults(settings Settings, err error) (Settings, error) {
	if Port != "" {
		settings.Port = Port
	} else {
		if settings.Port == "" {
			settings.Port = DEFAULTS["port"].(string)
		}
	}

	if Address != "" {
		settings.Address = Address
	} else {
		if settings.Address == "" {
			settings.Address = DEFAULTS["address"].(string)
		}
	}

	if settings.Scheme == "" {
		settings.Scheme = DEFAULTS["scheme"].(string)
	}

	if Workers > 0 {
		settings.Workers = Workers
	} else {
		settings.Workers = DEFAULTS["workers"].(int)
	}

	return settings, err
}

func BuildMain(configPath string) (Settings, error) {
	if _settings, err := build(Settings{}, configPath); err != nil {
		return _settings.(Settings), err
	} else {
		settings := _settings.(Settings)

		return settings, nil
	}
}

func BuildS3(configPath string) (S3, error) {
	if _settings, err := build(S3{}, configPath); err != nil {
		return _settings.(S3), err
	} else {
		settings := _settings.(S3)

		return settings, nil
	}
}

func build(result interface{}, path string) (interface{}, error) {
	if path == "" {
		return result, nil
	}

	var (
		err  error
		data []byte
	)

	path, err = filepath.Abs(path)
	if err != nil {
		return result, fmt.Errorf("Configuration: Path: %s", err)
	}

	data, err = ioutil.ReadFile(path)
	if err != nil {
		return result, fmt.Errorf("Configuration: Read: %s", err)
	}

	settingsArray := reflect.MakeMap(reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(result)))

	if err = yaml.Unmarshal(data, settingsArray.Interface()); err != nil {
		return result, fmt.Errorf("Configuration: Parse: %s", err)
	}

	if settings := settingsArray.MapIndex(reflect.ValueOf(Env)); settings.IsValid() {
		return settings.Interface(), nil
	} else {
		return result, fmt.Errorf("Not found config for env: %s", Env)
	}
}
