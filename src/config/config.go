package config

import "os"

type Config struct {
	ListenAddress string
}

func LoadConfig() Config {
	return Config{
		ListenAddress: getEnv("LISTEN_ADDRESS", ":8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
