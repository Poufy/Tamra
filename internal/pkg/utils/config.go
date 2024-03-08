package utils

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port                  int
	DBConn                string
	LogLevel              string
	FirebaseConfig        string
	RestaurantLogosBucket string
}

func GetConfig() Config {
	cfg := Config{}

	// Set up the flags with default values. If the environment variables are set, they will override these values. The third parameter is the default value.
	flag.IntVar(&cfg.Port, "port", getEnvAsInt("PORT", 8080), "Port for the application")
	flag.StringVar(&cfg.DBConn, "db", getEnv("DB_CONNECTION_STRING", "postgres://postgres:mysecretpassword@localhost:5432/tamra-postgis?sslmode=disable"), "Database connection string")
	flag.StringVar(&cfg.LogLevel, "log-level", getEnv("LOG_LEVEL", "debug"), "Log level")
	flag.StringVar(&cfg.FirebaseConfig, "firebase-config", getEnv("FIREBASE_CONFIG", "firebaseConfig.json"), "Path to the firebase config file")
	flag.StringVar(&cfg.RestaurantLogosBucket, "restaurant-logos-bucket", getEnv("RESTAURANT_LOGOS_BUCKET", "dev-tamra-restaurant-logos"), "Name of the bucket where restaurant logos are stored")
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
