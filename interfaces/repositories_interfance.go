package interfaces

import (
	"time"

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
	GetEmployeeByEmpID(empID string) (*models.Employees, error)
	GetAllEmployees(page int8, limit int8) ([]models.Employees, error)

	GetTotalEmployees() (int64, error)
	GetLatestEmployeeID() (int64, error)
}

type AttendanceRepository interface {
	CreateAttendance(attendance *models.Attendance) error
	GetAttendanceListToday(page int8, limit int8) ([]models.Attendance, error)
	GetAttendanceByDate(startDate time.Time, endDate time.Time, page int8, limit int8) ([]models.Attendance, error)
	GetTotalAttendanceByDate(startDate time.Time, endDate time.Time) (int64, error)
	GetTotalAttendanceToday() (int64, error)
	IsUserAttendToday(employeeID uint) (bool, error)
	GetAttendanceByDateRange(startDate time.Time, endDate time.Time, page int8, limit int8) ([]models.Attendance, error)
	GetTotalEmployeesToday() (int64, error)
	GetPresentEmployeesToday() ([]models.Attendance, error)

	// New methods for analytics
	GetAllAttendanceByDateRange(startDate time.Time, endDate time.Time) ([]models.Attendance, error)
	GetDailyAttendanceStats(startDate time.Time, endDate time.Time) ([]models.Attendance, error)
	GetUniqueEmployeesInDateRange(startDate time.Time, endDate time.Time) ([]models.Employees, error)
}
