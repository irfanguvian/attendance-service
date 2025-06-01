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

func (as *AttendanceService) GetTodayAttendanceSummary() (dto.TodayAttendanceSummary, error) {
	var result dto.TodayAttendanceSummary

	// Get today's date
	today := time.Now()
	result.Date = today.Format("2006-01-02")

	// Get total employees count
	totalEmployees, err := as.Repositories.AttendanceRepository.GetTotalEmployeesToday()
	if err != nil {
		return result, err
	}

	// Get attendance records for today
	attendances, err := as.Repositories.AttendanceRepository.GetPresentEmployeesToday()
	if err != nil {
		return result, err
	}

	// Count present employees using checkTime logic (similar to GetSalariesEmployeeByDate)
	var presentCount int64
	for _, attendance := range attendances {
		if checkTime(attendance.ClockIn) {
			presentCount++
		}
	}

	// Calculate absent count
	absentCount := totalEmployees - presentCount

	// Calculate rates
	var attendanceRate, absenceRate float64
	if totalEmployees > 0 {
		attendanceRate = (float64(presentCount) / float64(totalEmployees)) * 100
		absenceRate = (float64(absentCount) / float64(totalEmployees)) * 100
	}

	result.TotalEmployees = totalEmployees
	result.PresentCount = presentCount
	result.AbsentCount = absentCount
	result.AttendanceRate = attendanceRate
	result.AbsenceRate = absenceRate

	return result, nil
}

// Helper function to determine if an employee was late (present but after 9 AM)
func isLate(clockIn time.Time) bool {
	y, m, d := clockIn.Date()
	nineAM := time.Date(y, m, d, 9, 0, 0, 0, clockIn.Location())
	weekDayNum := clockIn.Weekday()
	isWeekday := weekDayNum >= time.Monday && weekDayNum <= time.Friday

	return isWeekday && clockIn.After(nineAM)
}

// Helper function to calculate working days between two dates
func calculateWorkingDays(startDate, endDate time.Time) int {
	count := 0
	for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 0, 1) {
		weekday := d.Weekday()
		if weekday >= time.Monday && weekday <= time.Friday {
			count++
		}
	}
	return count
}

// Helper function to categorize employee performance
func categorizePerformance(attendanceRate, onTimeRate float64) string {
	if attendanceRate >= 95 && onTimeRate >= 90 {
		return "Excellent"
	} else if attendanceRate >= 85 && onTimeRate >= 80 {
		return "Good"
	} else if attendanceRate >= 70 && onTimeRate >= 70 {
		return "Average"
	}
	return "Poor"
}

// GetComprehensiveAnalytics provides detailed analytics based on groupBy parameter
func (as *AttendanceService) GetComprehensiveAnalytics(startDate time.Time, endDate time.Time, groupBy string) (dto.ComprehensiveAttendanceAnalytics, error) {
	var result dto.ComprehensiveAttendanceAnalytics

	// Set date range
	result.DateRange = fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// Get all attendance data for the period
	attendances, err := as.Repositories.AttendanceRepository.GetAllAttendanceByDateRange(startDate, endDate)
	if err != nil {
		return result, err
	}

	// Get all unique employees in the date range
	employees, err := as.Repositories.AttendanceRepository.GetUniqueEmployeesInDateRange(startDate, endDate)
	if err != nil {
		return result, err
	}

	// Calculate working days in the period
	workingDays := calculateWorkingDays(startDate, endDate)

	// Build comprehensive analytics based on groupBy
	switch groupBy {
	case "daily":
		dailyMetrics, err := as.GetDailyTrendAnalytics(startDate, endDate)
		if err != nil {
			return result, err
		}
		result.DailyMetrics = dailyMetrics

	case "monthly":
		monthlyMetrics, err := as.GetMonthlyTrendAnalytics(startDate, endDate)
		if err != nil {
			return result, err
		}
		result.MonthlyMetrics = monthlyMetrics

	case "employee":
		employeeInsights := as.calculateEmployeeInsights(attendances, employees, workingDays)
		result.EmployeeInsights = employeeInsights

		// Get top and bottom performers
		result.TopPerformers = as.getTopPerformers(employeeInsights, 5)
		result.AttentionRequired = as.getBottomPerformers(employeeInsights, 5)

	default:
		// Return all analytics if no specific groupBy is provided
		dailyMetrics, _ := as.GetDailyTrendAnalytics(startDate, endDate)
		monthlyMetrics, _ := as.GetMonthlyTrendAnalytics(startDate, endDate)
		employeeInsights := as.calculateEmployeeInsights(attendances, employees, workingDays)

		result.DailyMetrics = dailyMetrics
		result.MonthlyMetrics = monthlyMetrics
		result.EmployeeInsights = employeeInsights
		result.TopPerformers = as.getTopPerformers(employeeInsights, 5)
		result.AttentionRequired = as.getBottomPerformers(employeeInsights, 5)
	}

	// Calculate summary
	result.Summary = as.calculateOverallSummary(attendances, employees, workingDays)

	// Calculate total salary disbursed
	result.TotalSalaryDisbursed = as.calculateTotalSalary(result.EmployeeInsights)

	// Generate insights
	result.Insights = as.generateInsights(result)

	return result, nil
}

// GetDailyTrendAnalytics provides daily attendance trends
func (as *AttendanceService) GetDailyTrendAnalytics(startDate time.Time, endDate time.Time) ([]dto.DailyAttendanceMetrics, error) {
	attendances, err := as.Repositories.AttendanceRepository.GetDailyAttendanceStats(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get total employees
	totalEmployees, err := as.Repositories.EmployeeRepository.GetTotalEmployees()
	if err != nil {
		return nil, err
	}

	// Group attendance by date
	dailyStats := make(map[string]map[string]int64)

	for _, attendance := range attendances {
		date := attendance.ClockIn.Format("2006-01-02")

		if dailyStats[date] == nil {
			dailyStats[date] = make(map[string]int64)
		}

		if checkTime(attendance.ClockIn) {
			dailyStats[date]["present"] = dailyStats[date]["present"] + 1
		} else if isLate(attendance.ClockIn) {
			dailyStats[date]["absent"] = dailyStats[date]["absent"] + 1
		}
	}

	var dailyMetrics []dto.DailyAttendanceMetrics

	// Generate metrics for each working day
	for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 0, 1) {
		weekday := d.Weekday()
		if weekday >= time.Monday && weekday <= time.Friday {
			date := d.Format("2006-01-02")
			stats := dailyStats[date]

			totalPresent := stats["present"]
			TotalAbsent := stats["absent"]

			var attendanceRate, onTimeRate, productivityScore float64
			if totalEmployees > 0 {
				attendanceRate = (float64(totalPresent) / float64(totalEmployees)) * 100
			}

			// Custom productivity score (attendance rate * on-time rate factor)
			productivityScore = attendanceRate*0.8 + onTimeRate*0.2

			dailyMetrics = append(dailyMetrics, dto.DailyAttendanceMetrics{
				Date:              date,
				TotalEmployees:    totalEmployees,
				PresentCount:      totalPresent,
				AbsentCount:       TotalAbsent,
				AttendanceRate:    attendanceRate,
				ProductivityScore: productivityScore,
			})
		}
	}

	return dailyMetrics, nil
}

// GetMonthlyTrendAnalytics provides monthly attendance and salary trends
func (as *AttendanceService) GetMonthlyTrendAnalytics(startDate time.Time, endDate time.Time) ([]dto.MonthlyAttendanceMetrics, error) {
	attendances, err := as.Repositories.AttendanceRepository.GetAllAttendanceByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get total employees
	totalEmployees, err := as.Repositories.EmployeeRepository.GetTotalEmployees()
	if err != nil {
		return nil, err
	}

	// Group data by month
	monthlyStats := make(map[string]map[string]interface{})

	for _, attendance := range attendances {
		monthKey := attendance.ClockIn.Format("2006-01")

		if monthlyStats[monthKey] == nil {
			monthlyStats[monthKey] = map[string]interface{}{
				"presentCount":     int64(0),
				"lateCount":        int64(0),
				"workingDays":      0,
				"employeePresence": make(map[string]map[string]int64),
			}
		}

		stats := monthlyStats[monthKey]
		empPresence := stats["employeePresence"].(map[string]map[string]int64)

		if empPresence[attendance.Employee.EmpID] == nil {
			empPresence[attendance.Employee.EmpID] = map[string]int64{
				"present": 0,
				"late":    0,
				"salary":  0,
			}
		}

		if checkTime(attendance.ClockIn) {
			stats["presentCount"] = stats["presentCount"].(int64) + 1
			empPresence[attendance.Employee.EmpID]["present"]++
		} else if isLate(attendance.ClockIn) {
			stats["lateCount"] = stats["lateCount"].(int64) + 1
			empPresence[attendance.Employee.EmpID]["late"]++
		}
	}

	var monthlyMetrics []dto.MonthlyAttendanceMetrics

	for monthKey, stats := range monthlyStats {
		// Parse month and year
		monthTime, _ := time.Parse("2006-01", monthKey)
		year := monthTime.Year()
		month := monthTime.Format("2006-01")

		// Calculate working days for the month
		monthStart := time.Date(year, monthTime.Month(), 1, 0, 0, 0, 0, monthTime.Location())
		monthEnd := monthStart.AddDate(0, 1, -1)

		if monthStart.Before(startDate) {
			monthStart = startDate
		}
		if monthEnd.After(endDate) {
			monthEnd = endDate
		}

		workingDays := calculateWorkingDays(monthStart, monthEnd)

		var avgAttendanceRate, avgOnTimeRate float64
		totalSalary := float32(0)

		empPresence := stats["employeePresence"].(map[string]map[string]int64)
		employeeCount := len(empPresence)

		var bestEmployee, worstEmployee string
		bestRate, worstRate := 0.0, 100.0

		for empID, presence := range empPresence {
			empPresent := presence["present"]
			empLate := presence["late"]
			empTotal := empPresent + empLate

			var empAttendanceRate, empOnTimeRate float64
			if workingDays > 0 {
				empAttendanceRate = (float64(empTotal) / float64(workingDays)) * 100
			}
			if empTotal > 0 {
				empOnTimeRate = (float64(empPresent) / float64(empTotal)) * 100
			}

			avgAttendanceRate += empAttendanceRate
			avgOnTimeRate += empOnTimeRate

			// Calculate salary for this employee for this month
			empSalary := (float32(empPresent) / 22.0) * 10000000.0
			totalSalary += empSalary

			// Track best and worst performers
			if empAttendanceRate > bestRate {
				bestRate = empAttendanceRate
				bestEmployee = empID
			}
			if empAttendanceRate < worstRate {
				worstRate = empAttendanceRate
				worstEmployee = empID
			}
		}

		if employeeCount > 0 {
			avgAttendanceRate = avgAttendanceRate / float64(employeeCount)
		}

		monthlyMetrics = append(monthlyMetrics, dto.MonthlyAttendanceMetrics{
			Month:                month,
			Year:                 year,
			TotalWorkingDays:     workingDays,
			TotalSalaryDisbursed: totalSalary,
			AvgAttendanceRate:    avgAttendanceRate,
			TotalEmployees:       totalEmployees,
			MostPresentEmployee:  bestEmployee,
			LeastPresentEmployee: worstEmployee,
		})
	}

	return monthlyMetrics, nil
}

// Helper methods for calculating insights
func (as *AttendanceService) calculateEmployeeInsights(attendances []models.Attendance, employees []models.Employees, workingDays int) []dto.EmployeeAttendanceInsight {
	employeeStats := make(map[string]map[string]int64)

	// Initialize stats for all employees
	for _, emp := range employees {
		employeeStats[emp.EmpID] = map[string]int64{
			"present": 0,
			"late":    0,
			"absent":  0,
		}
	}

	// Count attendance for each employee
	for _, attendance := range attendances {
		stats := employeeStats[attendance.Employee.EmpID]
		if checkTime(attendance.ClockIn) {
			stats["present"]++
		} else if isLate(attendance.ClockIn) {
			stats["late"]++
		} else {
			stats["absent"]++
		}
	}

	var insights []dto.EmployeeAttendanceInsight

	for _, emp := range employees {
		stats := employeeStats[emp.EmpID]
		present := stats["present"]
		late := stats["late"]
		totalPresent := present + late
		absent := int64(workingDays) - totalPresent

		var attendanceRate, onTimeRate float64
		if workingDays > 0 {
			attendanceRate = (float64(totalPresent) / float64(workingDays)) * 100
		}
		if totalPresent > 0 {
			onTimeRate = (float64(present) / float64(totalPresent)) * 100
		}

		// Calculate salary
		salary := (float32(present) / 22.0) * 10000000.0

		// Categorize performance
		category := categorizePerformance(attendanceRate, onTimeRate)

		insights = append(insights, dto.EmployeeAttendanceInsight{
			EmpID:               emp.EmpID,
			Fullname:            emp.Fullname,
			TotalPresent:        present,
			TotalAbsent:         absent,
			TotalLate:           late,
			AttendanceRate:      attendanceRate,
			OnTimeRate:          onTimeRate,
			TotalSalary:         salary,
			PerformanceCategory: category,
		})
	}

	return insights
}

func (as *AttendanceService) calculateOverallSummary(attendances []models.Attendance, employees []models.Employees, workingDays int) dto.AttendanceSummaryResponse {
	totalEmployees := int64(len(employees))
	presentCount := int64(0)

	for _, attendance := range attendances {
		if checkTime(attendance.ClockIn) || isLate(attendance.ClockIn) {
			presentCount++
		}
	}

	absentCount := (totalEmployees * int64(workingDays)) - presentCount

	var attendanceRate, absenceRate float64
	if totalEmployees > 0 && workingDays > 0 {
		attendanceRate = (float64(presentCount) / float64(totalEmployees*int64(workingDays))) * 100
		absenceRate = 100 - attendanceRate
	}

	return dto.AttendanceSummaryResponse{
		Date:           time.Now().Format("2006-01-02"),
		TotalEmployees: totalEmployees,
		PresentCount:   presentCount,
		AbsentCount:    absentCount,
		AttendanceRate: attendanceRate,
		AbsenceRate:    absenceRate,
	}
}

func (as *AttendanceService) getTopPerformers(insights []dto.EmployeeAttendanceInsight, count int) []dto.EmployeeAttendanceInsight {
	// Sort by attendance rate descending
	sorted := make([]dto.EmployeeAttendanceInsight, len(insights))
	copy(sorted, insights)

	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].AttendanceRate < sorted[j].AttendanceRate {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	if len(sorted) > count {
		return sorted[:count]
	}
	return sorted
}

func (as *AttendanceService) getBottomPerformers(insights []dto.EmployeeAttendanceInsight, count int) []dto.EmployeeAttendanceInsight {
	// Sort by attendance rate ascending
	sorted := make([]dto.EmployeeAttendanceInsight, len(insights))
	copy(sorted, insights)

	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].AttendanceRate > sorted[j].AttendanceRate {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	if len(sorted) > count {
		return sorted[:count]
	}
	return sorted
}

func (as *AttendanceService) calculateTotalSalary(insights []dto.EmployeeAttendanceInsight) float32 {
	total := float32(0)
	for _, insight := range insights {
		total += insight.TotalSalary
	}
	return total
}

func (as *AttendanceService) generateInsights(analytics dto.ComprehensiveAttendanceAnalytics) []string {
	var insights []string

	// Overall attendance insights
	if analytics.Summary.AttendanceRate >= 90 {
		insights = append(insights, "🎉 Excellent overall attendance rate! The team is highly engaged.")
	} else if analytics.Summary.AttendanceRate >= 80 {
		insights = append(insights, "👍 Good attendance rate, but there's room for improvement.")
	} else {
		insights = append(insights, "⚠️ Attendance rate needs attention. Consider reviewing policies.")
	}

	// Daily trend insights
	if len(analytics.DailyMetrics) > 0 {
		avgProductivity := 0.0
		for _, daily := range analytics.DailyMetrics {
			avgProductivity += daily.ProductivityScore
		}
		avgProductivity = avgProductivity / float64(len(analytics.DailyMetrics))

		if avgProductivity >= 80 {
			insights = append(insights, "📈 Productivity scores are consistently high across the period.")
		} else {
			insights = append(insights, "📉 Productivity scores indicate potential for improvement.")
		}
	}

	// Employee performance insights
	if len(analytics.EmployeeInsights) > 0 {
		excellentCount := 0
		poorCount := 0

		for _, emp := range analytics.EmployeeInsights {
			if emp.PerformanceCategory == "Excellent" {
				excellentCount++
			} else if emp.PerformanceCategory == "Poor" {
				poorCount++
			}
		}

		excellentPercent := (float64(excellentCount) / float64(len(analytics.EmployeeInsights))) * 100
		poorPercent := (float64(poorCount) / float64(len(analytics.EmployeeInsights))) * 100

		if excellentPercent >= 50 {
			insights = append(insights, fmt.Sprintf("🌟 %.0f%% of employees are excellent performers!", excellentPercent))
		}

		if poorPercent >= 20 {
			insights = append(insights, fmt.Sprintf("⚠️ %.0f%% of employees need performance support.", poorPercent))
		}
	}

	// Salary insights
	if analytics.TotalSalaryDisbursed > 0 {
		insights = append(insights, fmt.Sprintf("💰 Total salary disbursed: Rp %.0f", analytics.TotalSalaryDisbursed))
	}

	return insights
}
