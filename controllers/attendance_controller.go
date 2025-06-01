package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/irfanguvian/attendance-service/dto"
	"github.com/irfanguvian/attendance-service/interfaces"
	"github.com/irfanguvian/attendance-service/utils"
)

type AttendanceController struct {
	AttendanceService interfaces.AttendanceService
}

func NewAttendanceController(attendanceService interfaces.AttendanceService) *AttendanceController {
	return &AttendanceController{
		AttendanceService: attendanceService,
	}
}

func (ac *AttendanceController) CreateAttendance(c *gin.Context) {
	var reqBody dto.CreateAttendanceBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	err := ac.AttendanceService.CreateAttendance(reqBody)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, 200, "Attendance created successfully", nil)
}

func (ac *AttendanceController) GetAttendanceList(c *gin.Context) {
	var pagination dto.Pagination

	if _, ok := c.GetQuery("page"); ok {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid page number")
			return
		}
		pageUint := page
		pagination.Page = int8(pageUint)
	}

	if _, ok := c.GetQuery("limit"); ok {
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid limit number")
			return
		}
		limitUint := limit
		pagination.Limit = int8(limitUint)
	}

	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}

	attendances, err := ac.AttendanceService.GetAttendanceList(pagination)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Attendance list retrieved successfully", attendances)
}

func (ac *AttendanceController) ListEmployeeSalaries(c *gin.Context) {
	var pagination dto.PaginationEmployeeSalary

	if _, ok := c.GetQuery("page"); ok {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid page number")
			return
		}
		pageUint := page
		pagination.Page = int8(pageUint)
	}

	if _, ok := c.GetQuery("limit"); ok {
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid limit number")
			return
		}
		limitUint := limit
		pagination.Limit = int8(limitUint)
	}

	if _, ok := c.GetQuery("start_date"); ok {
		dateStr := c.Query("start_date")
		dateParse, err := time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid start date format")
			return
		}
		pagination.StartDate = dateParse
	} else {
		utils.ErrorResponse(c, 400, "Start Date is required")
		return
	}

	if _, ok := c.GetQuery("end_date"); ok {
		dateStr := c.Query("end_date")
		dateParse, err := time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid start date format")
			return
		}
		pagination.EndDate = dateParse
	} else {
		utils.ErrorResponse(c, 400, "End Date is required")
		return
	}

	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}

	paginationSend := dto.Pagination{
		Page:  pagination.Page,
		Limit: pagination.Limit,
	}

	salaries, err := ac.AttendanceService.GetSalariesEmployeeByDate(pagination.StartDate, pagination.EndDate, paginationSend)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, 200, "Employee salaries retrieved successfully", salaries)
}

func (ac *AttendanceController) GetAttendanceListByDateRange(c *gin.Context) {
	var pagination dto.AttendanceListWithDateRange

	if _, ok := c.GetQuery("page"); ok {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid page number")
			return
		}
		pagination.Page = int8(page)
	}

	if _, ok := c.GetQuery("limit"); ok {
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			utils.ErrorResponse(c, 400, "Invalid limit number")
			return
		}
		pagination.Limit = int8(limit)
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		utils.ErrorResponse(c, 400, "Both start_date and end_date are required (format: YYYY-MM-DD)")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid start_date format. Use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid end_date format. Use YYYY-MM-DD")
		return
	}

	// Ensure end date is end of day
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())

	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}

	paginationSend := dto.Pagination{
		Page:  pagination.Page,
		Limit: pagination.Limit,
	}

	attendances, err := ac.AttendanceService.GetAttendanceListByDateRange(startDate, endDate, paginationSend)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, 200, "Attendance list retrieved successfully", attendances)
}
