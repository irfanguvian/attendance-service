package interfaces

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/irfanguvian/attendance-service/dto"
	"github.com/irfanguvian/attendance-service/entities"
)

type Services struct {
	AuthService       AuthService
	EmployeeService   EmployeeService
	AttendanceService AttendanceService
}

type AuthService interface {
	Login(loginBody dto.LoginBody) (*dto.ResponseLoginService, error)
	Signup(signupBody dto.SignupBody) (string, error)
	SignOut(userID uint) error
	ExchangeToken(refreshToken string) (*dto.ResponseLoginService, error)

	ValidateToken(tokenString string) (jwt.MapClaims, error)

	GetUserByAccessID(accessID string) (*entities.User, error)

	GetAccessTokenByAccessID(accessID string) error
}

type EmployeeService interface {
	CreateEmployee(employee dto.CreateEmployeeBody) error
	UpdateEmployee(employee dto.UpdateEmployeeBody) error
	DeleteEmployee(employeeID uint) error
	GetEmployeeByID(employeeID uint) (*entities.Employees, error)
	GetAllEmployees(body dto.Pagination) (dto.ResponseGetAllEmployees, error)
}

type AttendanceService interface {
	CreateAttendance(attendance dto.CreateAttendanceBody) error
	GetAttendanceList(body dto.Pagination) (dto.ResponseGetAllAttendance, error)
	GetSalariesEmployeeByDate(startDate time.Time, endDate time.Time, body dto.Pagination) (dto.ResponseGetMetricSalariesByDate, error)
	GetAttendanceListByDateRange(startDate time.Time, endDate time.Time, body dto.Pagination) (dto.ResponseGetAllAttendance, error)
	GetTodayAttendanceSummary() (dto.TodayAttendanceSummary, error)

	// New analytics methods
	GetComprehensiveAnalytics(startDate time.Time, endDate time.Time, groupBy string) (dto.ComprehensiveAttendanceAnalytics, error)
	GetDailyTrendAnalytics(startDate time.Time, endDate time.Time) ([]dto.DailyAttendanceMetrics, error)
	GetMonthlyTrendAnalytics(startDate time.Time, endDate time.Time) ([]dto.MonthlyAttendanceMetrics, error)
}
