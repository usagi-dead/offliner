package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env" env-Default:"development"`
	DbPath     string     `yaml:"db_path" env:"DB_PATH" env-required:"true"`
	HttpServer HttpServer `yaml:"http_server"  env-required:"true"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

func MustLoad() *Config {
	ConfigPath := os.Getenv("CONFIG_PATH")
	if ConfigPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", ConfigPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(ConfigPath, &cfg); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
