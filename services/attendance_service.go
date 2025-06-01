package services

import (
	"errors"
	"math"

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
	checkAttend, err := as.Repositories.AttendanceRepository.IsUserAttendToday(attendance.EmployeeID)
	if err != nil {
		return err
	}

	if checkAttend {
		return errors.New("employee has already clocked in today")
	}
	newAttendance := &models.Attendance{
		EmployeeID: attendance.EmployeeID,
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
			ClockIn:    attendance.ClockIn,
			CreatedAt:  attendance.CreatedAt,
			UpdatedAt:  attendance.UpdatedAt,
		})
	}

	result.Total = getTotal
	result.TotalPage = math.Ceil(float64(getTotal) / float64(body.Limit))
	return result, nil
}
