package entities

import "time"

type Salary struct {
	ID          uint
	SalaryType  string
	TotalSalary int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
