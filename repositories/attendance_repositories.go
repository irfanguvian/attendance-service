package repositories

import (
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
	if err := ar.DB.Where("created_at >= CURRENT_DATE").Order("id desc").Offset(int(offset)).Limit(int(limit)).Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

func (ar *AttendanceRepositories) GetTotalAttendanceToday() (int64, error) {
	var count int64
	if err := ar.DB.Model(&models.Attendance{}).Where("created_at >= CURRENT_DATE").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (ar *AttendanceRepositories) IsUserAttendToday(employeeID uint) (bool, error) {
	var count int64
	if err := ar.DB.Model(&models.Attendance{}).
		Where("employee_id = ? AND created_at >= CURRENT_DATE", employeeID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}