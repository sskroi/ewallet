package config

import (
	"ewallet/pkg/storage/postgres"

	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PgCfg   postgres.Config `yaml:"postgresCfg" env-required:"true"`
	Port    string          `yaml:"port" env-required:"true"`
	GinMode string          `yaml:"ginMode" env-required:"true"`
}

func MustLoadCfg() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
