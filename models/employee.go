package models

import "time"

type Employees struct {
	ID        uint      `gorm:"primaryKey"`
	EmpID     string    `gorm:"unique;index:emp_id_idx"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Attandance *[]Attandance `gorm:"foreignKey:EmployeeID"`
}
