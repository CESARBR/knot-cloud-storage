package config

import (
	"os"
	"strings"

	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"github.com/spf13/viper"
)

// Server represents the server configuration properties
type Server struct {
	Port int
}

// Logger represents the logger configuration properties
type Logger struct {
	Level  string
	Syslog bool
}

// RabbitMQ represents the rabbitmq configuration properties
type RabbitMQ struct {
	URL string
}

// MongoDB represents the database configuration properties
type MongoDB struct {
	Host string
	Port int
	Name string
}

// Things represents the things service configuration properties
type Things struct {
	Host string
	Port int
}

// Expiration represents the data TTL configuration
type Expiration struct {
	Time int32
}

// Config represents the service configuration
type Config struct {
	Server
	Logger
	RabbitMQ
	MongoDB
	Things
	Expiration
}

func readFile(name string) {
	logger := logging.NewLogrus("error", false).Get("Config")
	viper.SetConfigName(name)
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatalf("Error reading config file, %s", err)
	}
}

// Load returns the service configuration
func Load() Config {
	var configuration Config
	logger := logging.NewLogrus("error", false).Get("Config")
	viper.AddConfigPath("internal/config")
	viper.SetConfigType("yaml")

	readFile("default")

	if os.Getenv("ENV") == "development" {
		readFile("development")
		if err := viper.MergeInConfig(); err != nil {
			logger.Fatalf("Error reading config file, %s", err)
		}
	}

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&configuration); err != nil {
		logger.Fatalf("Error unmarshalling configuration, %s", err)
	}

	return configuration
}
