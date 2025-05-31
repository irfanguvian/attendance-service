package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irfanguvian/attendance-service/config"
	"github.com/irfanguvian/attendance-service/controllers"
	"github.com/irfanguvian/attendance-service/interfaces"
	"github.com/irfanguvian/attendance-service/logger"
	"github.com/irfanguvian/attendance-service/middlewares"
	"github.com/irfanguvian/attendance-service/repositories"
	"github.com/irfanguvian/attendance-service/routes"
	"github.com/irfanguvian/attendance-service/services"
)

func init() {
	config.LoadConfig()

	config.ConnectDatabase()
}

func setupRepositories() interfaces.Repositories {
	return interfaces.Repositories{
		UserRepository: repositories.NewAuthRepositories(config.DB),
	}
}

func setupServices(repositories interfaces.Repositories) interfaces.Services {
	return interfaces.Services{
		AuthService: services.NewAuthService(repositories),
	}
}

func main() {
	router := gin.Default()

	// Health check
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	repo := setupRepositories()
	services := setupServices(repo)

	authController := controllers.NewAuthController(services.AuthService)

	controllers := &controllers.Controllers{
		AuthController: authController,
	}

	authMiddleware := middlewares.NewAuthMiddleware(services.AuthService)

	middlewareHandler := &middlewares.Middlewares{
		AuthMiddleware: authMiddleware,
	}

	routes.SetupRoutes(router, controllers, middlewareHandler)

	port := config.AppConfig.ServerPort

	logger.Info("Server starting on port %d", port)
	if err := router.Run(":" + strconv.Itoa(port)); err != nil {
		logger.Error("Failed to start server: %v", err)
		return
	}

}
