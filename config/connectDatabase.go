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
