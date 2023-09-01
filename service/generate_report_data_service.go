package service

import (
	"strings"
	"time"
	"work_report/model"
)

type GenerateReportDataService interface {
	Generate(works map[int]model.Work) model.Report
}

type generateReportDataService struct {
	model.Setting
}

func NewGenerateReportDataService(setting model.Setting) GenerateReportDataService {
	return &generateReportDataService{setting}
}

func (s generateReportDataService) Generate(works map[int]model.Work) model.Report {
	now := time.Now()
	year, month, _ := now.Date()
	targetMonth := time.Month(int(month) - s.Setting.MonthsAgo)

	lastOfCurrentMonth := time.Date(year, targetMonth+1, 0, 0, 0, 0, 0, time.Local)

	daysInCurrentMonth := lastOfCurrentMonth.Day()

	daily := model.DayToDailyData{}
	for date := 1; date <= daysInCurrentMonth; date++ {
		day := time.Date(year, targetMonth, date, 0, 0, 0, 0, time.Local)
		work := works[date-1]
		if work.Hours == 0 {
			continue
		}

		dailyData := model.DailyData{
			TargetDate: day.Format("2006-01-02"),
		}

		if works != nil {
			dailyReport := s.Setting.DailyReport
			dailyData.StartTime = dailyReport.StartsAt
			dailyData.RelaxTime = dailyReport.RestTime
			dailyData.WorkContent = strings.Join(work.Contents, ", ")
			_ = dailyData.SetWorkTime(work.HourMin())
		}

		daily[day.Format("20060102")] = dailyData
	}

	return model.Report{
		Data: model.WorkData{
			DailyReport: daily,
		},
	}
}
