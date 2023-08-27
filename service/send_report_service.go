package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"work_report/model"
)

type SendReportService interface {
	Send(report model.Report) error
}

type sendReportService struct {
	model.Setting
}

func NewSendReportService(setting model.Setting) SendReportService {
	return &sendReportService{setting}
}

func (s sendReportService) Send(report model.Report) error {
	url := fmt.Sprintf("https://platform.levtech.jp/p/workreport/input/%d/", s.Setting.ReportID)

	query, err := report.ToQuery()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", url)
	req.Header.Add("Cookie", "login=p")
	req.Header.Add("Cookie", fmt.Sprintf(
		"CAKEPHP=%s", s.Setting.SessionID,
	))
	req.Header.Add("Cookie", fmt.Sprintf(
		"AWSELBAuthSessionCookie-0=%s", s.Setting.AWSAuth,
	))

	client := http.DefaultClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Status Code: %d \n", resp.StatusCode)

	for k, v := range resp.Header {
		fmt.Printf("%s: %s\n", k, v)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf(string(body))
	}

	return nil
}
