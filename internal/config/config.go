package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// HTTPServer configures the HTTP server.
type HTTPServer struct {
	Addr string `yaml:"address"`
}

// Config holds all configuration for the application.
type Config struct {
	Env         string               `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string               `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"` // Embedded struct
}

// MustLoad reads the config file and returns the configuration. It panics on error.
func MustLoad() *Config {
	// Define the config path flag.
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to the configuration file (e.g., config/local.yaml)")
	flag.Parse()

	// If the flag is not set, try to get the path from an environment variable.
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	// If still no path, exit.
	if configPath == "" {
		log.Fatal("config path is not set: use -config flag or CONFIG_PATH env variable")
	}

	// Check if the file exists.
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	// Read the configuration file into the struct.
	// CRITICAL: We must pass a pointer to cfg using '&'.
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config file: %s", err)
	}

	return &cfg
}
