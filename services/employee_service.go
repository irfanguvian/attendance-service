package services

import (
	"fmt"
	"math"
	"time"

	"github.com/irfanguvian/attendance-service/dto"
	"github.com/irfanguvian/attendance-service/entities"
	"github.com/irfanguvian/attendance-service/interfaces"
	"github.com/irfanguvian/attendance-service/models"
)

type EmployeeService struct {
	Repositories interfaces.Repositories
}

func NewEmployeeService(repo interfaces.Repositories) interfaces.EmployeeService {
	return &EmployeeService{
		Repositories: repo,
	}
}

func GenerateUniqueID(prefix string, padding int, id int64) string {
	// Create the format string for zero-padding.
	// For example, if padding is 4, the format will be "%04d".
	format := fmt.Sprintf("%%0%dd", padding)

	// Format the integer with the created format string and combine it with the prefix.
	return fmt.Sprintf("%s-%s", prefix, fmt.Sprintf(format, id))
}

func (es *EmployeeService) CreateEmployee(employee dto.CreateEmployeeBody) error {
	getLatestEmployee, err := es.Repositories.EmployeeRepository.GetLatestEmployeeID()
	if err != nil {
		return err
	}

	empID := GenerateUniqueID("EMP", 5, getLatestEmployee+1)
	newEmployee := &models.Employees{
		Fullname: employee.Fullname,
		EmpID:    empID,
	}

	return es.Repositories.EmployeeRepository.CreateEmployee(newEmployee)
}

func (es *EmployeeService) UpdateEmployee(employee dto.UpdateEmployeeBody) error {
	existingEmployee, err := es.Repositories.EmployeeRepository.GetEmployeeByID(employee.EmployeeID)
	if err != nil {
		return err
	}

	existingEmployee.Fullname = employee.Fullname

	return es.Repositories.EmployeeRepository.UpdateEmployee(existingEmployee)
}

func (es *EmployeeService) DeleteEmployee(employeeID uint) error {
	existingEmployee, err := es.Repositories.EmployeeRepository.GetEmployeeByID(employeeID)
	if err != nil {
		return err
	}
	now := time.Now()

	existingEmployee.DeletedAt = &now

	return es.Repositories.EmployeeRepository.UpdateEmployee(existingEmployee)
}

func (es *EmployeeService) GetEmployeeByID(employeeID uint) (*entities.Employees, error) {
	existingEmployee, err := es.Repositories.EmployeeRepository.GetEmployeeByID(employeeID)
	if err != nil {
		return nil, err
	}

	return &entities.Employees{
		ID:       existingEmployee.ID,
		EmpID:    existingEmployee.EmpID,
		Fullname: existingEmployee.Fullname,
	}, nil
}

func (es *EmployeeService) GetAllEmployees(body dto.Pagination) (dto.ResponseGetAllEmployees, error) {
	var result dto.ResponseGetAllEmployees
	result.Limit = body.Limit
	result.Page = body.Page

	employees, err := es.Repositories.EmployeeRepository.GetAllEmployees(body.Page, body.Limit)
	if err != nil {
		return result, err
	}

	getTotal, err := es.Repositories.EmployeeRepository.GetTotalEmployees()
	if err != nil {
		return result, err
	}

	if len(employees) == 0 {
		return result, nil // Return nil if no employees found
	}

	for _, emp := range employees {
		Employees := entities.Employees{
			ID:        emp.ID,
			EmpID:     emp.EmpID,
			Fullname:  emp.Fullname,
			CreatedAt: emp.CreatedAt,
			UpdatedAt: emp.UpdatedAt,
		}

		result.Employees = append(result.Employees, Employees)
	}

	result.Total = getTotal
	result.TotalPage = math.Ceil(float64(getTotal) / float64(body.Limit))

	return result, nil
}
