package controllers

import (
	"strconv"

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
