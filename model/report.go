package model

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
)

type Report struct {
	Data WorkData `json:"data" url:"data"`
}

type WorkData struct {
	DailyReport DayToDailyData `json:"DailyReport" url:"DailyReport"`
}

type DayToDailyData map[string]DailyData

func (d DayToDailyData) EncodeValues(key string, val *url.Values) error {
	for day, dailyData := range d {
		if dailyData.ID != "" {
			val.Set(fmt.Sprintf("%s[%s][id]", key, day), dailyData.ID)
		}
		val.Set(fmt.Sprintf("%s[%s][target_date]", key, day), dailyData.TargetDate)
		val.Set(fmt.Sprintf("%s[%s][start_time]", key, day), dailyData.StartTime)
		val.Set(fmt.Sprintf("%s[%s][end_time]", key, day), dailyData.EndTime)
		val.Set(fmt.Sprintf("%s[%s][relax_time]", key, day), dailyData.RelaxTime)
		val.Set(fmt.Sprintf("%s[%s][work_content]", key, day), dailyData.WorkContent)
		val.Set(fmt.Sprintf("%s[%s][work_time]", key, day), dailyData.WorkTime)
	}

	return nil
}

type DailyData struct {
	ID          string `json:"id" url:"id"`
	TargetDate  string `json:"target_date" url:"target_date"`
	StartTime   string `json:"start_time" url:"start_time"`
	EndTime     string `json:"end_time" url:"end_time"`
	RelaxTime   string `json:"relax_time" url:"relax_time"`
	WorkTime    string `json:"work_time" url:"work_time"`
	WorkContent string `json:"work_content" url:"work_content"`
}

func FileToReport(filePath string) (*Report, error) {
	bin, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var report Report
	if err := json.Unmarshal(bin, &report); err != nil {
		return nil, err
	}

	return &report, nil
}

func (r Report) ToFile(filePath string) error {
	file, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, file, 0644); err != nil {
		return err
	}

	return nil
}

func (r Report) ToQuery() (string, error) {
	vals, err := query.Values(r)
	if err != nil {
		return "", err
	}

	return vals.Encode(), nil
}

func (d *DailyData) CalcWorkTime() error {
	startedHourMin := strings.Split(d.StartTime, ":")
	startedHour, err := strconv.Atoi(startedHourMin[0])
	if err != nil {
		return err
	}
	startedMin, err := strconv.Atoi(startedHourMin[1])
	if err != nil {
		return err
	}
	finishedHourMin := strings.Split(d.EndTime, ":")
	finishedHour, err := strconv.Atoi(finishedHourMin[0])
	if err != nil {
		return err
	}
	finishedMin, err := strconv.Atoi(finishedHourMin[1])
	if err != nil {
		return err
	}
	restHourMin := strings.Split(d.RelaxTime, ":")
	restHour, err := strconv.Atoi(restHourMin[0])
	if err != nil {
		return err
	}
	restMin, err := strconv.Atoi(restHourMin[1])
	if err != nil {
		return err
	}

	totalWorkedMinutes :=
		(finishedHour-startedHour-restHour)*60 +
			(finishedMin - startedMin - restMin)
	workedHour := totalWorkedMinutes / 60
	workedMin := totalWorkedMinutes % 60
	d.WorkTime = fmt.Sprintf("%02d:%02d", workedHour, workedMin)
	return nil
}
