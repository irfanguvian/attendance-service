package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// code to connect database
	var err error
	dsn := AppConfig.DatabaseURL
	
	// Log configuration for debugging (redact sensitive info in production)
	fmt.Println("Database URL:", maskDatabaseURL(AppConfig.DatabaseURL))
	fmt.Println("Environment:", AppConfig.Environment)
	
	if dsn == "" {
		log.Fatal("Database URL is not configured. Please check your environment variables.")
	}
	
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	
	// Verify connection
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error getting database connection: %v", err)
	}
	
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database connection verification failed: %v", err)
	}
	
	fmt.Println("Database connected successfully")
}

// Helper function to mask sensitive information in database URL for logging
func maskDatabaseURL(url string) string {
	if len(url) < 20 {
		return "[DATABASE_URL_NOT_SET_OR_INVALID]"
	}
	// Return just enough information to identify the database type and host
	return url[:15] + "..." + url[len(url)-10:]
}
