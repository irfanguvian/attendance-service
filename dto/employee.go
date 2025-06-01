package dto

import "github.com/irfanguvian/attendance-service/entities"

type CreateEmployeeBody struct {
	Fullname string `json:"fullname" binding:"required"`
}

type UpdateEmployeeBody struct {
	EmployeeID uint   `json:"employee_id" binding:"required"`
	Fullname   string `json:"fullname" binding:"required"`
}

type PaginationEmployeeSalary struct {
	Pagination
	DateFilter string `json:"date_filter" binding:"required"`
}

type ResponseEntitySalaryData struct {
	entities.Employees
	Salary int32 `json:"salary"`
}

type ResponseGetAllEmployees struct {
	Employees []entities.Employees `json:"employees"`
	PaginationResponse
}