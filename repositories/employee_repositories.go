package repositories

import (
	"time"

	"github.com/irfanguvian/attendance-service/models"
	"gorm.io/gorm"
)

type EmployeeRepositories struct {
	DB *gorm.DB
}

func NewEmployeeRepositories(db *gorm.DB) *EmployeeRepositories {
	return &EmployeeRepositories{
		DB: db,
	}
}

func (er *EmployeeRepositories) CreateEmployee(employee *models.Employees) error {
	if err := er.DB.Create(employee).Error; err != nil {
		return err
	}
	return nil
}

func (er *EmployeeRepositories) GetEmployeeByID(employeeID uint) (*models.Employees, error) {
	employee := &models.Employees{}
	if err := er.DB.Where("deleted_at is null").First(employee, employeeID).Error; err != nil {
		return nil, err
	}
	return employee, nil
}

func (er *EmployeeRepositories) UpdateEmployee(employee *models.Employees) error {
	if err := er.DB.Save(employee).Error; err != nil {
		return err
	}
	return nil
}

func (er *EmployeeRepositories) DeleteEmployee(employeeID uint) error {
	//soft delete employee
	if err := er.DB.Where("id = ?", employeeID).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

func (er *EmployeeRepositories) GetAllEmployees(page int8, limit int8) ([]models.Employees, error) {
	var employees []models.Employees
	offset := (page - 1) * limit
	if err := er.DB.Where("deleted_at is null").Order("id desc").Offset(int(offset)).Limit(int(limit)).Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}

func (er *EmployeeRepositories) GetTotalEmployees() (int64, error) {
	var count int64
	if err := er.DB.Model(&models.Employees{}).Where("deleted_at is null").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (er *EmployeeRepositories) GetLatestEmployeeID() (int64, error) {
	var latestEmployee models.Employees
	if err := er.DB.Order("id desc").First(&latestEmployee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil // No employees found
		}
		return 0, err
	}
	return int64(latestEmployee.ID), nil
}
