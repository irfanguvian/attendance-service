package interfaces

import (
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
	ListEmployeeSalaries(body dto.PaginationEmployeeSalary) ([]dto.ResponseEntitySalaryData, error)
}

type AttendanceService interface {
	CreateAttendance(attendance dto.CreateAttendanceBody) error
	GetAttendanceList(body dto.Pagination) (dto.ResponseGetAllAttendance, error)
	// GetMetricAttendanceByDate(date string) (*dto.ResponseGetMetricAttendanceByDate, error)
}

// type SalariesService interface {
// 	CreateSalary(salary dto.CreateSalaryBody) (*dto.ResponseCreateSalary, error)
// 	UpdateSalary(salary dto.UpdateSalaryBody) (*dto.ResponseUpdateSalary, error)
// 	GetMetricSalariesByDate(date string) (*dto.ResponseGetMetricSalariesByDate, error)
// }
