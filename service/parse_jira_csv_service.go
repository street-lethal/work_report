package service

import (
	"strconv"
	"time"
	"work_report/model"
)

type ParseJiraCSVService interface {
	Parse(csv [][]string) (map[int]model.Work, error)
}

type parseJiraCSVService struct {
}

func NewParseJiraCSVService() ParseJiraCSVService {
	return &parseJiraCSVService{}
}

func (ss parseJiraCSVService) Parse(csv [][]string) (map[int]model.Work, error) {
	works := map[int]model.Work{}
	for i, row := range csv {
		if i == 0 { // skip header row
			continue
		}

		ticket := row[0]
		hour, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, err
		}

		date, err := time.Parse("2006-01-02 15:04", row[4])
		if err != nil {
			return nil, err
		}
		day := date.Day()

		work, ok := works[day]
		if !ok {
			work = model.Work{}
		}
		work.AddContent(ticket)
		work.AddHour(hour)
		works[day] = work
	}

	return works, nil
}
