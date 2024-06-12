package configs

import (
	"fmt"
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

func getMySqlConnectionString() string {
	userName := getEnv("MYSQL_ACC", "")
	password := getEnv("MYSQL_PASSWORD", "")
	host := getEnv("MYSQL_HOST", "")
	port := getEnv("MYSQL_PORT", "")
	dbName := getEnv("MYSQL_DB_NAME", "")
	parseTime := "true"
	charset := "utf8mb4"
	location := "Local"
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=%s&charset=%s&loc=%s",
		userName, password, host, port, dbName, parseTime, charset, location,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

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
		MySqlConnectionString:  getMySqlConnectionString(),
		JwtSecret:              []byte(getEnv("JWT_SECRET", "")),
		EnableHttpsMode:        getEnv("HTTPS_MODE", "") == "true",
	}
}
