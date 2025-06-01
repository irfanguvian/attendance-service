package dto

import (
	"time"

	"github.com/irfanguvian/attendance-service/entities"
)

type CreateAttendanceBody struct {
	EmployeeID uint      `json:"employee_id" binding:"required"`
	ClockIn    time.Time `json:"clock_in" binding:"required"` 
}

type ResponseGetAllAttendance struct {
	Attendance []entities.Attendance `json:"attendance"`
	PaginationResponse
}
