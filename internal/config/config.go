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

		HotelHost string `env:"HOTEL_HOST" envDefault:"localhost"`
		HotelPort string `env:"HOTEL_PORT" envDefault:"50052"`

		BookingHost string `env:"BOOKING_HOST" envDefault:"localhost"`
		BookingPort string `env:"BOOKING_PORT" envDefault:"50053"`
	}

	Log struct {
		LogLvl string `env:"LOGGER_LEVEL" envDefault:"info"`
	}
)

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("configs file path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("configs file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("configs path is empty: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to configs file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
