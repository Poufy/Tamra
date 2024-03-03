package utils

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port   int
	DBConn string
}

func GetConfig() Config {
	cfg := Config{}

	// Set up the flags with default values. If the environment variables are set, they will override these values. The third parameter is the default value.
	flag.IntVar(&cfg.Port, "port", getEnvAsInt("PORT", 8080), "Port for the application")
	flag.StringVar(&cfg.DBConn, "db", getEnv("DB_CONNECTION_STRING", "postgres://user:pass@localhost/dbname"), "Database connection string")

	flag.Parse()

	fmt.Printf("Configuration values: %v\n", cfg)
	return cfg
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultVal
}
