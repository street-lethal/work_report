package model

import (
	"encoding/json"
	"os"
)

type Setting struct {
	MonthsAgo   int `json:"months_ago"`
	DailyReport struct {
		StartsAt string `json:"starts_at"`
		RestTime string `json:"rest_time"`
	} `json:"daily_report"`
	ReportID int `json:"report_id"`
}

func FileToSetting(filePath string) (*Setting, error) {
	bin, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var setting Setting
	if err := json.Unmarshal(bin, &setting); err != nil {
		return nil, err
	}

	return &setting, nil
}
