package models

import "time"

type Employees struct {
	ID        uint      `gorm:"primaryKey"`
	EmpID     string    `gorm:"unique;index:emp_id_idx"`
	Fullname  string    `gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"default:null" json:"delete_at,omitempty"`

	Attandance *[]Attendance `gorm:"foreignKey:EmployeeID"`
}
