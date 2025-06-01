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
	Absent   int8
	Present  int8
	EmpID    string
	Fullname string
}

type ResponseGetMetricSalariesByDate struct {
	Salary []ResponseEntitySalaryData `json:"salary"`
	PaginationResponse
}

// New DTOs for attendance summary and analytics
type AttendanceSummaryResponse struct {
	Date           string  `json:"date"`
	TotalEmployees int64   `json:"total_employees"`
	PresentCount   int64   `json:"present_count"`
	AbsentCount    int64   `json:"absent_count"`
	AttendanceRate float64 `json:"attendance_rate"`
	AbsenceRate    float64 `json:"absence_rate"`
}

type DailyTrendData struct {
	Date           string  `json:"date"`
	PresentCount   int64   `json:"present_count"`
	AbsentCount    int64   `json:"absent_count"`
	AttendanceRate float64 `json:"attendance_rate"`
}

type MonthlyTrendData struct {
	Month             string  `json:"month"`
	Year              int     `json:"year"`
	TotalSalary       float32 `json:"total_salary"`
	AvgAttendanceRate float64 `json:"avg_attendance_rate"`
}

type AttendanceAnalyticsResponse struct {
	Summary      AttendanceSummaryResponse `json:"summary"`
	DailyTrend   []DailyTrendData          `json:"daily_trend"`
	MonthlyTrend []MonthlyTrendData        `json:"monthly_trend"`
}

type DateRangeRequest struct {
	StartDate string `json:"start_date" form:"start_date"`
	EndDate   string `json:"end_date" form:"end_date"`
}

type AttendanceListWithDateRange struct {
	Pagination
	DateRangeRequest
}
