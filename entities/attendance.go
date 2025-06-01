package entities

import "time"

type Attendance struct {
	ID         uint
	EmployeeID uint
	ClockIn    time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
