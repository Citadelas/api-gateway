package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env      string   `yaml:"env" env-default:"local"`
	Services Services `yaml:"services"`
}

type Services struct {
	Task Service `yaml:"task"`
	SSO  Service `yaml:"sso"`
}

type Service struct {
	Endpoint string        `yaml:"endpoint"`
	timeout  time.Duration `yaml:"timeout"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file doesn't exists" + path)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config " + err.Error())
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
