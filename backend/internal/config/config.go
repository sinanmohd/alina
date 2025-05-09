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
	CorsAllowAll  bool   `toml:"cors_allow_all"`
}
type DatabaseConfig struct {
	Url string `toml:"url"`
}
type Config struct {
	Server ServerConfig   `toml:"server"`
	Db     DatabaseConfig `toml:"db"`
}

func New() (*Config, error) {
	defaultSecretKey := "change-me-for-dev-only"
	defaultDataDir := "alina_data"
	if value, ok := os.LookupEnv("STATE_DIRECTORY"); ok {
		defaultDataDir = value
	}

	config := Config{
		Server: ServerConfig{
			Host:          "[::]",
			Port:          8008,
			Data:          defaultDataDir,
			PublicUrl:     "http://localhost:8008",
			FileSizeLimit: 1024 * 1024 * 256, // 256MB
			ChunkSize:     1024 * 1024, // 1MB
			SecretKey:     defaultSecretKey,
			CorsAllowAll:  false,
		},
		Db: DatabaseConfig{
			Url: "postgresql:///alina?user=alina&host=/var/run/postgresql",
		},
	}

	var configPath string
	defaultConfigPath := "/etc/alina.conf"
	if value, ok := os.LookupEnv("ALINA_CONFIG"); ok {
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

	flag.StringVar(&config.Server.Host, "host", config.Server.Host, "Bind host for alina")
	flag.UintVar(&config.Server.Port, "port", config.Server.Port, "Bind port for alina")
	flag.StringVar(&config.Server.PublicUrl, "public-url", config.Server.PublicUrl, "Public URL for alina")
	flag.StringVar(&config.Server.Data, "data", config.Server.Data, "Data directory for alina")
	flag.BoolVar(&config.Server.CorsAllowAll, "cors-allow-all", config.Server.CorsAllowAll, "Allow requsts from all cors origins")
	flag.StringVar(&config.Db.Url, "db-url", config.Db.Url, "Database URL")
	flag.Parse()

	if value, ok := os.LookupEnv("DB_URL"); ok {
		config.Db.Url = value
	}

	if value, ok := os.LookupEnv("ALINA_SECRET_KEY"); ok {
		config.Server.SecretKey = value
	}
	if (config.Server.SecretKey == defaultSecretKey) {
		log.Println("Warning using the defaultSecretKey is not safe for production")
	}

	// TODO: validator

	return &config, nil
}
