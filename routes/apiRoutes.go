package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/irfanguvian/attendance-service/controllers"
	"github.com/irfanguvian/attendance-service/middlewares"
)

func SetupRoutes(router *gin.Engine, controllers *controllers.Controllers, middleware *middlewares.Middlewares) {
	authRoutes := router.Group("/auth")
	authRoutes.POST("/login", controllers.AuthController.Login)
	authRoutes.POST("/signup", controllers.AuthController.Signup)
	authRoutes.POST("/signout", middleware.AuthMiddleware.ProtecHandlerRequest, controllers.AuthController.SignOut)
	authRoutes.POST("/exchange-token", controllers.AuthController.ExchangeToken)
}
