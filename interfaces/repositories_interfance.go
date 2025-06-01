package interfaces

import (
	"github.com/irfanguvian/attendance-service/models"
)

type Repositories struct {
	UserRepository       UserRepository
	EmployeeRepository   EmployeeRepository
	AttendanceRepository AttendanceRepository
}

type UserRepository interface {
	CreateUser(email, password string) (uint, error)
	CreateAccessToken(userID uint, accessID string) error
	CreateRefreshToken(accessID string, refreshID string) error

	DeleteAccessTokenByUserID(userID uint) error

	GetUserByEmail(email string) (*models.User, error)
	GetRefreshTokenByAccessID(accessID string) (string, error)

	GetUserByAccessID(accessID string) (*models.User, error)
}

type EmployeeRepository interface {
	CreateEmployee(employee *models.Employees) error
	UpdateEmployee(employee *models.Employees) error
	DeleteEmployee(employeeID uint) error
	GetEmployeeByID(employeeID uint) (*models.Employees, error)
	GetAllEmployees(page int8, limit int8) ([]models.Employees, error)

	GetTotalEmployees() (int64, error)
	GetLatestEmployeeID() (int64, error)
}

type AttendanceRepository interface {
	CreateAttendance(attendance *models.Attendance) error
	GetAttendanceListToday(page int8, limit int8) ([]models.Attendance, error)
	GetTotalAttendanceToday() (int64, error)
	IsUserAttendToday(employeeID uint) (bool, error)
}
