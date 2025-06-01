package dto

import (
	"time"

	"github.com/irfanguvian/attendance-service/entities"
)

type CreateAttendanceBody struct {
	EmpID   string    `json:"emp_id" binding:"required"`
	ClockIn time.Time `json:"clock_in" binding:"required"`
}

type ResponseGetAllAttendance struct {
	Attendance []entities.Attendance `json:"attendance"`
	PaginationResponse
}

type SalaryCounterObject struct {
	Absent  int8
	Present int8
	EmpID  string
	Fullname string
}

type ResponseGetMetricSalariesByDate struct {
	Salary []ResponseEntitySalaryData `json:"salary"`
	PaginationResponse
}
