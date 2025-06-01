package services

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/irfanguvian/attendance-service/dto"
	"github.com/irfanguvian/attendance-service/entities"
	"github.com/irfanguvian/attendance-service/interfaces"
	"github.com/irfanguvian/attendance-service/models"
)

type AttendanceService struct {
	Repositories interfaces.Repositories
}

func NewAttendanceService(repo interfaces.Repositories) interfaces.AttendanceService {
	return &AttendanceService{
		Repositories: repo,
	}
}

func (as *AttendanceService) CreateAttendance(attendance dto.CreateAttendanceBody) error {

	employee, err := as.Repositories.EmployeeRepository.GetEmployeeByEmpID(attendance.EmpID)

	if err != nil {
		return err
	}

	checkAttend, err := as.Repositories.AttendanceRepository.IsUserAttendToday(employee.ID)
	if err != nil {
		return err
	}

	if checkAttend {
		return errors.New("employee already attended today")
	}

	newAttendance := &models.Attendance{
		EmployeeID: employee.ID,
		ClockIn:    attendance.ClockIn,
	}

	return as.Repositories.AttendanceRepository.CreateAttendance(newAttendance)
}

func (as *AttendanceService) GetAttendanceList(body dto.Pagination) (dto.ResponseGetAllAttendance, error) {
	var result dto.ResponseGetAllAttendance
	result.Limit = body.Limit
	result.Page = body.Page

	attendances, err := as.Repositories.AttendanceRepository.GetAttendanceListToday(body.Page, body.Limit)
	if err != nil {
		return result, err
	}
	getTotal, err := as.Repositories.AttendanceRepository.GetTotalAttendanceToday()
	if err != nil {
		return result, err
	}

	for _, attendance := range attendances {
		result.Attendance = append(result.Attendance, entities.Attendance{
			ID:         attendance.ID,
			EmployeeID: attendance.EmployeeID,
			Employee: entities.Employees{
				EmpID:    attendance.Employee.EmpID,
				Fullname: attendance.Employee.Fullname,
			},
			ClockIn:   attendance.ClockIn,
			CreatedAt: attendance.CreatedAt,
			UpdatedAt: attendance.UpdatedAt,
		})
	}

	result.Total = getTotal
	result.TotalPage = math.Ceil(float64(getTotal) / float64(body.Limit))
	return result, nil
}

func checkTime(checkTime time.Time) bool {
	// Extract the year, month, and day from the time to be checked.
	y, m, d := checkTime.Date()

	// Create the 9:00 AM threshold for that specific date and location.
	// Using the original time's location is crucial for correctness.
	nineAM := time.Date(y, m, d, 9, 0, 0, 0, checkTime.Location())

	isBeforeOrEqualNineAM := !checkTime.After(nineAM)
	weekDayNum := checkTime.Weekday()
	isWeekday := weekDayNum >= time.Monday && weekDayNum <= time.Friday

	return isBeforeOrEqualNineAM && isWeekday

}

func (ss *AttendanceService) GetSalariesEmployeeByDate(startDate time.Time, endDate time.Time, body dto.Pagination) (dto.ResponseGetMetricSalariesByDate, error) {
	var result dto.ResponseGetMetricSalariesByDate
	result.Limit = body.Limit
	result.Page = body.Page

	getAttendance, err := ss.Repositories.AttendanceRepository.GetAttendanceByDate(startDate, endDate, body.Page, body.Limit)
	if err != nil {
		return result, err
	}
	totalAttendanceRow, err := ss.Repositories.AttendanceRepository.GetTotalAttendanceByDate(startDate, endDate)
	if err != nil {
		return result, err
	}
	// count salary for each employee
	salariesCounter := make(map[string]dto.SalaryCounterObject)

	for _, attendance := range getAttendance {
		isAbsent := false
		if !checkTime(attendance.ClockIn) {
			isAbsent = true
		}

		SalaryCheck, ok := salariesCounter[attendance.Employee.EmpID]

		if !ok {
			if isAbsent {
				salariesCounter[attendance.Employee.EmpID] = dto.SalaryCounterObject{
					Absent:   1,
					Present:  0,
					EmpID:    attendance.Employee.EmpID,
					Fullname: attendance.Employee.Fullname,
				}
			} else {
				salariesCounter[attendance.Employee.EmpID] = dto.SalaryCounterObject{
					Absent:   0,
					Present:  1,
					EmpID:    attendance.Employee.EmpID,
					Fullname: attendance.Employee.Fullname,
				}
			}
		} else {
			if isAbsent {
				SalaryCheck.Absent++
			} else {
				SalaryCheck.Present++
			}
			salariesCounter[attendance.Employee.EmpID] = SalaryCheck
		}
	}

	fmt.Println("salariesCounter", salariesCounter["EMP-00001"])

	var salaryData []dto.ResponseEntitySalaryData

	for _, val := range salariesCounter {
		countSalary := (float32(val.Present) / 22.0) * 10000000.0

		salary := dto.ResponseEntitySalaryData{
			EmpID:           val.EmpID,
			Fullname:        val.Fullname,
			Salary:          countSalary,
			TotalAttendance: int8(val.Present + val.Absent),
			Absent:          val.Absent,
			Present:         val.Present,
		}

		salaryData = append(salaryData, salary)

	}
	// Calculate total pages
	totalPages := math.Ceil(float64(totalAttendanceRow) / float64(body.Limit))
	result.Total = totalAttendanceRow
	result.TotalPage = totalPages
	result.Salary = salaryData

	return result, nil
}
func (as *AttendanceService) GetAttendanceListByDateRange(startDate time.Time, endDate time.Time, body dto.Pagination) (dto.ResponseGetAllAttendance, error) {
	var result dto.ResponseGetAllAttendance
	result.Limit = body.Limit
	result.Page = body.Page

	attendances, err := as.Repositories.AttendanceRepository.GetAttendanceByDateRange(startDate, endDate, body.Page, body.Limit)
	if err != nil {
		return result, err
	}

	getTotal, err := as.Repositories.AttendanceRepository.GetTotalAttendanceByDate(startDate, endDate)
	if err != nil {
		return result, err
	}

	for _, attendance := range attendances {
		result.Attendance = append(result.Attendance, entities.Attendance{
			ID:         attendance.ID,
			EmployeeID: attendance.EmployeeID,
			Employee: entities.Employees{
				EmpID:    attendance.Employee.EmpID,
				Fullname: attendance.Employee.Fullname,
			},
			ClockIn:   attendance.ClockIn,
			CreatedAt: attendance.CreatedAt,
			UpdatedAt: attendance.UpdatedAt,
		})
	}

	result.Total = getTotal
	result.TotalPage = math.Ceil(float64(getTotal) / float64(body.Limit))

	return result, nil
}
