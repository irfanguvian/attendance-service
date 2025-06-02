package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName              string
	ServerPort           int
	DatabaseURL          string
	JWTSecretKey         string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

var AppConfig Config

func LoadConfig() {
	// Load .env file
	env, err := godotenv.Read()
	isUsingProduction := false
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
		isUsingProduction = true
	}

	// Set defaults
	AppConfig = Config{
		AppName:              "Attendance System",
		ServerPort:           3000,
		AccessTokenDuration:  10 * time.Minute,
		RefreshTokenDuration: 1 * time.Hour,
	}

	if isUsingProduction {
		if port, err := strconv.Atoi(os.Getenv("PORT")); err == nil && port > 0 {
			AppConfig.ServerPort = port
		}

		if dbURL := os.Getenv("DATA_BASE_URL"); dbURL != "" {
			AppConfig.DatabaseURL = dbURL
		}

		if jwtKey := os.Getenv("JWT_SECRET_KEY"); jwtKey != "" {
			AppConfig.JWTSecretKey = jwtKey
		} else {
			AppConfig.JWTSecretKey = "secret" // Default for development only
			log.Println("Warning: Using default JWT secret key. This should be set in production.")
		}

	} else {
		if port, err := strconv.Atoi(os.Getenv("PORT")); err == nil && port > 0 {
			AppConfig.ServerPort = port
		}

		if dbURL := env["DATA_BASE_URL"]; dbURL != "" {
			AppConfig.DatabaseURL = dbURL
		}

		if jwtKey := env["JWT_SECRET_KEY"]; jwtKey != "" {
			AppConfig.JWTSecretKey = jwtKey
		} else {
			AppConfig.JWTSecretKey = "secret" // Default for development only
			log.Println("Warning: Using default JWT secret key. This should be set in production.")
		}

	}
}
