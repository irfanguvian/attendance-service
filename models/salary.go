package models

import "time"

type Salary struct {
	ID          uint `gorm:"primaryKey"`
	SalaryType  string `json:"salary_type"`
	TotalSalary int32 `json:"total_salary"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
