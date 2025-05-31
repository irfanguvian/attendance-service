package models

import "time"

type Attandance struct {
	ID         uint      `gorm:"primaryKey"`
	EmployeeID uint      `gorm:"index:employee_id_idx"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Employee Employees `gorm:"foreignKey:EmployeeID"`
}
