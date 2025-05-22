package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type (
	Config struct {
		App    App
		Logger Log
		GW     Gateway
	}

	App struct {
		ServiceName string `env:"SERVICE_NAME" envDefault:"Unnamed_Service"`
		Version     string `env:"VERSION" envDefault:"1.0.0"`
	}

	Gateway struct {
		Port string `env:"HTTP_PORT" envDefault:"8080"`
		Host string `env:"HTTP_HOST" envDefault:"localhost"`

		AuthHost string `env:"AUTH_HOST" envDefault:"localhost"`
		AuthPort string `env:"AUTH_PORT" envDefault:"50051"`
	}

	Log struct {
		LogLvl string `env:"LOGGER_LEVEL" envDefault:"info"`
	}
)

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config file path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
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
