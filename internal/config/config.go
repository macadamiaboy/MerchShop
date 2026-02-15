package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

//may fall because of auth part in yaml

type Config struct {
	Server struct {
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		Timeout     int    `yaml:"timeout"`
		IdleTimeout int    `yaml:"idle_timout"`
	} `yaml:"server"`

	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`

	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
}

func LoadDBConfigData() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("failed to load data")
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("cannot read config %s", err)
	}

	return &cfg
}
