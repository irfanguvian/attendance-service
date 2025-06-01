package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irfanguvian/attendance-service/dto"
	"github.com/irfanguvian/attendance-service/interfaces"
	"github.com/irfanguvian/attendance-service/utils"
)

type EmployeeController struct {
	EmployeeService interfaces.EmployeeService
}

func NewEmployeeController(employeeService interfaces.EmployeeService) *EmployeeController {
	return &EmployeeController{
		EmployeeService: employeeService,
	}
}
func (ec *EmployeeController) CreateEmployee(c *gin.Context) {
	var reqBody dto.CreateEmployeeBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	err := ec.EmployeeService.CreateEmployee(reqBody)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, 200, "Employee created successfully", nil)
}

func (ec *EmployeeController) UpdateEmployee(c *gin.Context) {
	var reqBody dto.UpdateEmployeeBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	err := ec.EmployeeService.UpdateEmployee(reqBody)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, 200, "Employee updated successfully", nil)
}

func (ec *EmployeeController) DeleteEmployee(c *gin.Context) {
	employeeID := c.Param("employeeID")
	if employeeID == "" {
		utils.ErrorResponse(c, 400, "Employee ID is required")
		return
	}
	// Convert employeeID to uint
	employeeIDUint, err := strconv.ParseUint(employeeID, 10, 64)

	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid Employee ID format")
		return
	}

	err = ec.EmployeeService.DeleteEmployee(uint(employeeIDUint))
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, 200, "Employee deleted successfully", nil)
}

func (ec *EmployeeController) GetEmployeeByID(c *gin.Context) {
	employeeID := c.Param("employeeID")
	if employeeID == "" {
		utils.ErrorResponse(c, 400, "Employee ID is required")
		return
	}
	// Convert employeeID to uint
	employeeIDUint, err := strconv.ParseUint(employeeID, 10, 64)

	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid Employee ID format")
		return
	}

	employee, err := ec.EmployeeService.GetEmployeeByID(uint(employeeIDUint))
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, 200, "Employee retrieved successfully", employee)
}

func (ec *EmployeeController) GetAllEmployees(c *gin.Context) {
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

	employees, err := ec.EmployeeService.GetAllEmployees(pagination)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, 200, "Employees retrieved successfully", employees)
}

func (ec *EmployeeController) ListEmployeeSalaries(c *gin.Context) {
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

	if _, ok := c.GetQuery("date_filter"); ok {
		pagination.DateFilter = c.Query("date_filter")
	} else {
		utils.ErrorResponse(c, 400, "Date filter is required")
		return
	}

	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}

	salaries, err := ec.EmployeeService.ListEmployeeSalaries(pagination)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, 200, "Employee salaries retrieved successfully", salaries)
}
