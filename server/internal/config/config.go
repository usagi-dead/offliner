package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env              string           `yaml:"env" env-Default:"development"`
	DbConfig         DbConfig         `yaml:"db"`
	HttpServerConfig HttpServerConfig `yaml:"http_server"  env-required:"true"`
	CacheConfig      CacheConfig      `yaml:"cache"`
	SMTPConfig       SMTPConfig       `yaml:"smtp"`
	JWTConfig        JWTConfig        `yaml:"jwt"`
}

type CacheConfig struct {
	Address                      string        `yaml:"address" env-required:"true"`
	Db                           int           `yaml:"db"`
	StateExpiration              time.Duration `yaml:"state_expiration" env-required:"true"`
	EmailConfirmedCodeExpiration time.Duration `yaml:"email_confirmed_code_expiration" env-required:"true"`
}

type HttpServerConfig struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

type DbConfig struct {
	Username string `yaml:"username"`
	Address  string `yaml:"address"`
	DbName   string `yaml:"db_name"`
	Sslmode  string `yaml:"sslmode"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
}

type JWTConfig struct {
	SigningMethod string        `yaml:"signing_method" env-required:"true"`
	AccessExpire  time.Duration `yaml:"access_expire" env-required:"true"`
	RefreshExpire time.Duration `yaml:"refresh_expire" env-required:"true"`
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
