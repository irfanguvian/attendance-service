package repositories

import (
	"time"

	"github.com/irfanguvian/attendance-service/models"
	"gorm.io/gorm"
)

type AttendanceRepositories struct {
	DB *gorm.DB
}

func NewAttendanceRepositories(db *gorm.DB) *AttendanceRepositories {
	return &AttendanceRepositories{
		DB: db,
	}
}

func (ar *AttendanceRepositories) CreateAttendance(attendance *models.Attendance) error {
	if err := ar.DB.Create(attendance).Error; err != nil {
		return err
	}
	return nil
}

func (ar *AttendanceRepositories) GetAttendanceListToday(page int8, limit int8) ([]models.Attendance, error) {
	var attendances []models.Attendance
	offset := (page - 1) * limit
	if err := ar.DB.Joins("Employee").Where("attendances.clock_in >= CURRENT_DATE").Order("id desc").Offset(int(offset)).Limit(int(limit)).Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

func (ar *AttendanceRepositories) GetTotalAttendanceToday() (int64, error) {
	var count int64
	if err := ar.DB.Model(&models.Attendance{}).Where("clock_in >= CURRENT_DATE").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (ar *AttendanceRepositories) IsUserAttendToday(employeeID uint) (bool, error) {
	var count int64
	if err := ar.DB.Model(&models.Attendance{}).
		Where("employee_id = ? AND clock_in >= CURRENT_DATE", employeeID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (ar *AttendanceRepositories) GetAttendanceByDate(startDate time.Time, endDate time.Time, page int8, limit int8) ([]models.Attendance, error) {
	var attendances []models.Attendance
	offset := (page - 1) * limit
	if err := ar.DB.
		Preload("Employee").
		Joins("JOIN employees Employee ON attendances.employee_id = Employee.id"). // Adjust 'attendances.employee_id' based on your actual foreign key
		Select("attendances.*, Employee.*").
		Where("Employee.deleted_at is null AND attendances.clock_in BETWEEN ? AND ?", startDate, endDate).
		Order("attendances.id ASC").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

func (ar *AttendanceRepositories) GetTotalAttendanceByDate(startDate time.Time, endDate time.Time) (int64, error) {
	var count int64
	if err := ar.DB.Model(&models.Attendance{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (ar *AttendanceRepositories) GetAttendanceByDateRange(startDate time.Time, endDate time.Time, page int8, limit int8) ([]models.Attendance, error) {
	var attendances []models.Attendance
	offset := (page - 1) * limit
	if err := ar.DB.
		Preload("Employee").
		Joins("JOIN employees Employee ON attendances.employee_id = Employee.id").
		Select("attendances.*, Employee.*").
		Where("Employee.deleted_at is null AND attendances.clock_in BETWEEN ? AND ?", startDate, endDate).
		Order("attendances.created_at DESC").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

func (ar *AttendanceRepositories) GetTotalEmployeesToday() (int64, error) {
	var count int64
	if err := ar.DB.Model(&models.Employees{}).Where("deleted_at IS NULL").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (ar *AttendanceRepositories) GetPresentEmployeesToday() ([]models.Attendance, error) {
	var attendances []models.Attendance
	if err := ar.DB.
		Preload("Employee").
		Joins("JOIN employees Employee ON attendances.employee_id = Employee.id").
		Select("attendances.*, Employee.*").
		Where("Employee.deleted_at IS NULL AND attendances.clock_in >= CURRENT_DATE").
		Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

// New methods for analytics
func (ar *AttendanceRepositories) GetAllAttendanceByDateRange(startDate time.Time, endDate time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	if err := ar.DB.
		Preload("Employee").
		Joins("JOIN employees Employee ON attendances.employee_id = Employee.id").
		Select("attendances.*, Employee.*").
		Where("Employee.deleted_at IS NULL AND attendances.clock_in BETWEEN ? AND ?", startDate, endDate).
		Order("attendances.created_at ASC").
		Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

func (ar *AttendanceRepositories) GetDailyAttendanceStats(startDate time.Time, endDate time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	if err := ar.DB.
		Preload("Employee").
		Joins("JOIN employees Employee ON attendances.employee_id = Employee.id").
		Select("attendances.*, Employee.*").
		Where("Employee.deleted_at IS NULL AND attendances.clock_in BETWEEN ? AND ?", startDate, endDate).
		Order("attendances.created_at ASC").
		Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

func (ar *AttendanceRepositories) GetUniqueEmployeesInDateRange(startDate time.Time, endDate time.Time) ([]models.Employees, error) {
	var employees []models.Employees
	if err := ar.DB.
		Distinct().
		Joins("JOIN attendances ON employees.id = attendances.employee_id").
		Where("employees.deleted_at IS NULL AND attendances.clock_in BETWEEN ? AND ?", startDate, endDate).
		Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}
