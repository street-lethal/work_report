package service

import (
	"strings"
	"time"
	"work_report/model"
)

type GenerateReportDataService interface {
	Generate(works map[int][]string) model.Report
}

type generateReportDataService struct {
	model.Setting
}

func NewGenerateReportDataService(setting model.Setting) GenerateReportDataService {
	return &generateReportDataService{setting}
}

func (s generateReportDataService) Generate(works map[int][]string) model.Report {
	now := time.Now()
	year, month, _ := now.Date()
	targetMonth := time.Month(int(month) - s.Setting.MonthsAgo)

	lastOfCurrentMonth := time.Date(year, targetMonth+1, 0, 0, 0, 0, 0, time.Local)

	daysInCurrentMonth := lastOfCurrentMonth.Day()

	daily := model.DayToDailyData{}
	for date := 1; date <= daysInCurrentMonth; date++ {
		day := time.Date(year, targetMonth, date, 0, 0, 0, 0, time.Local)
		if day.Weekday() == time.Saturday || day.Weekday() == time.Sunday || s.isHoliday(date) {
			continue
		}

		dailyData := model.DailyData{
			TargetDate: day.Format("2006-01-02"),
		}

		if works != nil {
			dailyReport := s.Setting.DailyReport
			dailyData.StartTime = dailyReport.StartsAt
			dailyData.EndTime = dailyReport.EndsAt
			dailyData.RelaxTime = dailyReport.RestTime
			dailyData.WorkContent = strings.Join(works[date-1], ", ")
			_ = dailyData.CalcWorkTime()
		}

		daily[day.Format("20060102")] = dailyData
	}

	return model.Report{
		Data: model.WorkData{
			DailyReport: daily,
		},
	}
}

func (s generateReportDataService) isHoliday(date int) bool {
	for _, holiday := range s.Setting.Holidays {
		if date == holiday {
			return true
		}
	}

	return false
}
