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

// DTO for today's attendance summary
type TodayAttendanceSummary struct {
	Date           string  `json:"date"`
	TotalEmployees int64   `json:"total_employees"`
	PresentCount   int64   `json:"present_count"`
	AbsentCount    int64   `json:"absent_count"`
	AttendanceRate float64 `json:"attendance_rate"`
	AbsenceRate    float64 `json:"absence_rate"`
}

// Enhanced DTOs for comprehensive analytics
type DailyAttendanceMetrics struct {
	Date              string  `json:"date"`
	TotalEmployees    int64   `json:"total_employees"`
	PresentCount      int64   `json:"present_count"`
	AbsentCount       int64   `json:"absent_count"`
	AttendanceRate    float64 `json:"attendance_rate"`    // (Present / Total) * 100
	ProductivityScore float64 `json:"productivity_score"` // Custom score based on attendance patterns
}

type MonthlyAttendanceMetrics struct {
	Month                string  `json:"month"` // Format: "2024-01"
	Year                 int     `json:"year"`
	TotalWorkingDays     int     `json:"total_working_days"` // Working days in the month
	TotalSalaryDisbursed float32 `json:"total_salary_disbursed"`
	AvgAttendanceRate    float64 `json:"avg_attendance_rate"`
	TotalEmployees       int64   `json:"total_employees"`
	MostPresentEmployee  string  `json:"most_present_employee"`
	LeastPresentEmployee string  `json:"least_present_employee"`
}

type EmployeeAttendanceInsight struct {
	EmpID               string  `json:"emp_id"`
	Fullname            string  `json:"fullname"`
	TotalPresent        int64   `json:"total_present"`
	TotalAbsent         int64   `json:"total_absent"`
	TotalLate           int64   `json:"total_late"`
	AttendanceRate      float64 `json:"attendance_rate"`
	OnTimeRate          float64 `json:"on_time_rate"`
	TotalSalary         float32 `json:"total_salary"`
	PerformanceCategory string  `json:"performance_category"` // "Excellent", "Good", "Average", "Poor"
}

type AttendanceAnalyticsRequest struct {
	StartDate string `json:"start_date" form:"start_date" binding:"required"`
	EndDate   string `json:"end_date" form:"end_date" binding:"required"`
	GroupBy   string `json:"group_by" form:"group_by"` // "daily", "monthly", "employee"
}

type ComprehensiveAttendanceAnalytics struct {
	DateRange            string                      `json:"date_range"`
	Summary              AttendanceSummaryResponse   `json:"summary"`
	DailyMetrics         []DailyAttendanceMetrics    `json:"daily_metrics,omitempty"`
	MonthlyMetrics       []MonthlyAttendanceMetrics  `json:"monthly_metrics,omitempty"`
	EmployeeInsights     []EmployeeAttendanceInsight `json:"employee_insights,omitempty"`
	TopPerformers        []EmployeeAttendanceInsight `json:"top_performers,omitempty"`     // Top 5 performers
	AttentionRequired    []EmployeeAttendanceInsight `json:"attention_required,omitempty"` // Bottom 5 performers
	TotalSalaryDisbursed float32                     `json:"total_salary_disbursed"`
	Insights             []string                    `json:"insights"` // AI-like insights about patterns
}
