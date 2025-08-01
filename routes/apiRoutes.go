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

	employeeRoutes := router.Group("/employee")
	employeeRoutes.POST("/create", middleware.AuthMiddleware.ProtecHandlerRequest, controllers.EmployeeController.CreateEmployee)
	employeeRoutes.PUT("/update", middleware.AuthMiddleware.ProtecHandlerRequest, controllers.EmployeeController.UpdateEmployee)
	employeeRoutes.DELETE("/:employeeID", middleware.AuthMiddleware.ProtecHandlerRequest, controllers.EmployeeController.DeleteEmployee)
	employeeRoutes.GET("/:employeeID", middleware.AuthMiddleware.ProtecHandlerRequest, controllers.EmployeeController.GetEmployeeByID)
	employeeRoutes.GET("/list", middleware.AuthMiddleware.ProtecHandlerRequest, controllers.EmployeeController.GetAllEmployees)

	attendanceRoutes := router.Group("/attendance")
	attendanceRoutes.POST("/create", middleware.AuthMiddleware.ProtecHandlerRequest, controllers.AttendanceController.CreateAttendance)
	attendanceRoutes.GET("/list-today", middleware.AuthMiddleware.ProtecHandlerRequest, controllers.AttendanceController.GetAttendanceList)
	attendanceRoutes.GET("/list-by-date-range", middleware.AuthMiddleware.ProtecHandlerRequest, controllers.AttendanceController.GetAttendanceListByDateRange)
	attendanceRoutes.GET("/salaries", middleware.AuthMiddleware.ProtecHandlerRequest, controllers.AttendanceController.ListEmployeeSalaries)
	attendanceRoutes.GET("/summary-today", middleware.AuthMiddleware.ProtecHandlerRequest, controllers.AttendanceController.GetTodayAttendanceSummary)

	// Analytics routes
	attendanceRoutes.GET("/analytics/comprehensive", middleware.AuthMiddleware.ProtecHandlerRequest, controllers.AttendanceController.GetComprehensiveAnalytics)
	attendanceRoutes.GET("/analytics/daily-trends", middleware.AuthMiddleware.ProtecHandlerRequest, controllers.AttendanceController.GetDailyTrendAnalytics)
	attendanceRoutes.GET("/analytics/monthly-trends", middleware.AuthMiddleware.ProtecHandlerRequest, controllers.AttendanceController.GetMonthlyTrendAnalytics)

}
