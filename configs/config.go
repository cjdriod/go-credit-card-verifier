package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type ApiStatus struct {
	Success string
	Fail    string
}

type Config struct {
	ApiStatus
	EnablePremiumCardCheck bool
	MySqlConnectionString  string
	JwtSecret              []byte
	EnableHttpsMode        bool
}

var Constant = initConfig()

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		ApiStatus: ApiStatus{
			Success: "success",
			Fail:    "failed",
		},
		EnablePremiumCardCheck: getEnv("ENABLE_PREMIUM_CARD_CHECK", "") == "true",
		MySqlConnectionString:  getEnv("MYSQL_CONNECTION_STRING", ""),
		JwtSecret:              []byte(getEnv("JWT_SECRET", "")),
		EnableHttpsMode:        getEnv("HTTPS_MODE", "") == "true",
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
