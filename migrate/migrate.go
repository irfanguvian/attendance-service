package main

import (
	"fmt"
	"log"

	"github.com/irfanguvian/attendance-service/config"
	"github.com/irfanguvian/attendance-service/models"
)

func init() {
	config.LoadConfig()
	config.ConnectDatabase()
}

func main() {
	res := config.DB.Migrator().CurrentDatabase()
	fmt.Println(res)
	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("AutoMigrate User failed: %v", err)
	}
	if err := config.DB.AutoMigrate(&models.AccessToken{}); err != nil {
		log.Fatalf("AutoMigrate AccessToken failed: %v", err)
	}
	if err := config.DB.AutoMigrate(&models.RefreshToken{}); err != nil {
		log.Fatalf("AutoMigrate RefreshToken failed: %v", err)
	}
	if err := config.DB.AutoMigrate(&models.Salary{}); err != nil {
		log.Fatalf("AutoMigrate Salary failed: %v", err)
	}

	if err := config.DB.AutoMigrate(&models.Employees{}); err != nil {
		log.Fatalf("AutoMigrate Employees failed: %v", err)
	}

	if err := config.DB.AutoMigrate(&models.Attandance{}); err != nil {
		log.Fatalf("AutoMigrate Attandance failed: %v", err)
	}

	log.Println("Migration completed successfully")
}
