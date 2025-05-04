package config

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
	Host          string `toml:"host"`
	Port          uint   `toml:"port"`
	Data          string `toml:"data"`
	PublicUrl     string `toml:"public_url"`
	FileSizeLimit int    `toml:"file_size_limit"`
	ChunkSize     int    `toml:"chunk_size"`
	SecretKey     string `toml:"secret_key"`
}
type DatabaseConfig struct {
	Url string `toml:"url"`
}
type Config struct {
	Server ServerConfig   `toml:"server"`
	Db     DatabaseConfig `toml:"db"`
}

func New() (*Config, error) {
	var config Config = Config{
		Server: ServerConfig{
			Host:          "localhost",
			Port:          8008,
			Data:          "alina_data",
			PublicUrl:     "http://localhost",
			FileSizeLimit: 134217728,
			ChunkSize:     1048576,
			SecretKey:     "change-me-for-dev-only",
		},
		Db: DatabaseConfig{
			Url: "postgresql:///alina?user=alina&host=/var/run/postgresql",
		},
	}

	var configPath string
	defaultConfigPath := "/etc/alina.conf"
	if value, ok := os.LookupEnv("CONFIG"); ok {
		configPath = value
	} else {
		configPath = defaultConfigPath
	}

	_, err := os.Stat(configPath)
	if err != nil {
		if configPath == defaultConfigPath && !errors.Is(err, os.ErrNotExist) {
			log.Println("Error reading default config:", err)
			return nil, err
		} else if configPath != defaultConfigPath {
			log.Println("Error reading config:", err)
			return nil, err
		}
	}

	if _, err := os.Stat(configPath); err == nil {
		_, err := toml.DecodeFile(configPath, &config)
		if err != nil {
			log.Println("Error decoding TOML:", err)
			return nil, err
		}
	}

	flag.StringVar(&config.Server.Host, "server-host", config.Server.Host, "Bind host for alina")
	flag.UintVar(&config.Server.Port, "server-port", config.Server.Port, "Bind port for alina")
	flag.StringVar(&config.Server.PublicUrl, "server-public-url", config.Server.PublicUrl, "Public URL for alina")
	flag.StringVar(&config.Server.Data, "server-data", config.Server.Data, "Data directory for alina")
	flag.StringVar(&config.Db.Url, "db-url", config.Db.Url, "Database URL")
	flag.Parse()

	if value, ok := os.LookupEnv("DB_URL"); ok {
		config.Db.Url = value
	}

	// TODO: validator

	return &config, nil
}
