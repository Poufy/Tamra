package utils

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Port                  int
	DBConn                string
	LogLevel              string
	FirebaseConfigJSON    string
	RestaurantLogosBucket string
	Stage                 string
}

func GetConfig() Config {
	cfg := Config{}

	// Set up the flags with default values. If the environment variables are set, they will override these values. The third parameter is the default value.
	flag.IntVar(&cfg.Port, "port", getEnvAsInt("PORT", 8080), "Port for the application")
	flag.StringVar(&cfg.DBConn, "db", getEnv("DB_CONNECTION_STRING", "postgres://postgres:mysecretpassword@localhost:5432/tamra-postgis?sslmode=disable"), "Database connection string")
	flag.StringVar(&cfg.LogLevel, "log-level", getEnv("LOG_LEVEL", "debug"), "Log level")
	flag.StringVar(&cfg.FirebaseConfigJSON, "firebase-config-json", getEnv("FIREBASE_CONFIG_JSON", ""), "JSON string of the configuration for Firebase Authentication.")
	flag.StringVar(&cfg.RestaurantLogosBucket, "restaurant-logos-bucket", getEnv("RESTAURANT_LOGOS_BUCKET", "dev-tamra-restaurant-logos"), "Name of the bucket where restaurant logos are stored")
	flag.StringVar(&cfg.Stage, "stage", getEnv("STAGE", "dev"), "Stage of the application")
	flag.Parse()

	fmt.Printf("Configuration values: %v\n", cfg)
	logrus.Warn("firebase config json: ", cfg.FirebaseConfigJSON)
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
