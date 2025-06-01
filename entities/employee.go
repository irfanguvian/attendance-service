package entities

import "time"

type Employees struct {
	ID uint 
	EmpID string 
	Fullname string
	CreatedAt time.Time 
	UpdatedAt time.Time 
	DeletedAt *time.Time
}

