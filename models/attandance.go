package models

import "time"

type Attendance struct {
	ID         uint      `gorm:"primaryKey"`
	EmployeeID uint      `gorm:"index:employee_id_idx"`
	ClockIn    time.Time `json:"clock_in" gorm:"type:timestamp;not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Employee Employees `gorm:"foreignKey:EmployeeID"`
}
