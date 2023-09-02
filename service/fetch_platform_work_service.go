package service

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
	"work_report/model"

	"golang.org/x/net/html"
)

type FetchPlatformWorkService interface {
	FetchReports(session model.PlatformSession) (string, error)
	FindTargetMonthID(html *string) (string, error)
	FetchReport(reportID string, session model.PlatformSession) (string, error)
	FindDailyIDs(htm *string) (map[int]string, error)
	AttachDailyIDsToReport(
		ids map[int]string, report model.Report,
	) *model.Report
}

type fetchPlatformWorkService struct {
	setting model.Setting
	ParseHTMLService
}

func NewFetchPlatformWorkService(
	setting model.Setting,
	ps ParseHTMLService,
) FetchPlatformWorkService {
	return &fetchPlatformWorkService{setting, ps}
}

func (s fetchPlatformWorkService) FetchReports(
	session model.PlatformSession,
) (string, error) {
	fmt.Println("=== Fetch Reports ===")
	url := "https://platform.levtech.jp/p/workreport/"

	return s.fetch(url, session)
}

func (s fetchPlatformWorkService) FindTargetMonthID(htm *string) (string, error) {
	node, err := s.ParseHTMLService.Parse(*htm)
	if err != nil {
		return "", err
	}

	now := time.Now()
	targetMonth := now.AddDate(0, -s.setting.MonthsAgo, 0)
	strTargetMonth := targetMonth.Format("2006/01")
	node = s.ParseHTMLService.FindFirst(
		node, func(n *html.Node) bool {
			return n.Type == html.ElementNode &&
				n.Data == "a" &&
				n.FirstChild != nil &&
				n.FirstChild.Data == strTargetMonth
		},
	)

	link := s.ParseHTMLService.GetAttrValue(node, "href") // /p/workreport/12345/

	re := regexp.MustCompile(`\bworkreport/\d+`) // workreport/12345
	sliced := string(re.Find([]byte(link)))
	re = regexp.MustCompile(`\d+`) // 12345
	id := string(re.Find([]byte(sliced)))
	if id == "" {
		return "", fmt.Errorf("cannot take id from workreport path")
	}

	return id, nil
}

func (s fetchPlatformWorkService) FetchReport(
	reportID string, session model.PlatformSession,
) (string, error) {
	fmt.Println("=== Fetch Report ===")
	url := fmt.Sprintf("https://platform.levtech.jp/p/workreport/input/%s/", reportID)

	return s.fetch(url, session)
}

func (s fetchPlatformWorkService) FindDailyIDs(htm *string) (map[int]string, error) {
	node, err := s.ParseHTMLService.Parse(*htm)
	if err != nil {
		return nil, err
	}

	dailyIDs := map[int]string{}

	now := time.Now()
	targetMonth := now.AddDate(0, -s.setting.MonthsAgo, 0)
	for day := 1; day <= 31; day++ {
		date := fmt.Sprintf(
			"%s%02d",
			targetMonth.Format("200601"), day,
		)
		node := s.ParseHTMLService.FindFirst(
			node, func(n *html.Node) bool {
				if n.Type != html.ElementNode || n.Data != "input" {
					return false
				}

				return s.ParseHTMLService.GetAttrValue(n, "name") ==
					fmt.Sprintf("data[DailyReport][%s][id]", date)
			},
		)
		if node == nil {
			continue
		}

		id := s.ParseHTMLService.GetAttrValue(node, "value")
		if id != "" {
			dailyIDs[day] = id
		}
	}

	return dailyIDs, nil
}

func (s fetchPlatformWorkService) AttachDailyIDsToReport(
	ids map[int]string, report model.Report,
) *model.Report {
	year, month, _ := time.Now().AddDate(0, -s.setting.MonthsAgo, 0).Date()
	for day, id := range ids {
		reportKey := fmt.Sprintf("%04d%02d%02d", year, month, day)
		daily := report.Data.DailyReport[reportKey]
		daily.ID = id
		report.Data.DailyReport[reportKey] = daily
	}

	return &report
}

func (s fetchPlatformWorkService) fetch(
	url string, session model.PlatformSession,
) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.AddCookie(&http.Cookie{Name: "login", Value: "p"})
	req.AddCookie(&http.Cookie{Name: "CAKEPHP", Value: session.SessionID})
	req.AddCookie(&http.Cookie{
		Name: "AWSELBAuthSessionCookie-0", Value: session.AWSAuth,
	})

	client := http.DefaultClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Printf("Status Code: %d \n", resp.StatusCode)

	for k, v := range resp.Header {
		fmt.Printf("%s: %s\n", k, v)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf(string(body))
	}

	return string(body), nil
}
