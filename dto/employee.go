package dto

import (
	"time"

	"github.com/irfanguvian/attendance-service/entities"
)

type CreateEmployeeBody struct {
	Fullname string `json:"fullname" binding:"required"`
}

type UpdateEmployeeBody struct {
	EmployeeID uint   `json:"employee_id" binding:"required"`
	Fullname   string `json:"fullname" binding:"required"`
}

type PaginationEmployeeSalary struct {
	Pagination
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type ResponseEntitySalaryData struct {
	EmpID           string `json:"emp_id"`
	Fullname        string `json:"fullname"`
	Salary          float32  `json:"salary"`
	TotalAttendance int8   `json:"total_attendance"`
	Absent          int8   `json:"absent"`
	Present         int8   `json:"present"`
}

type ResponseGetAllEmployees struct {
	Employees []entities.Employees `json:"employees"`
	PaginationResponse
}
