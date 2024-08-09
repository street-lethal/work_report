package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"work_report/model"
)

type LoginPlatformService interface {
	Login(id model.PlatformID) (*model.PlatformSession, error)
}

type loginPlatformService struct {
}

func NewLoginPlatformService() LoginPlatformService {
	return &loginPlatformService{}
}

func (s loginPlatformService) Login(id model.PlatformID) (*model.PlatformSession, error) {
	client := http.DefaultClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	url := "https://platform.levtech.jp/p/"
	url, cookies, err := s.req1(client, url)
	if err != nil {
		return nil, err
	}

	pathWithQuery, cookies, err := s.req2(client, url, cookies)
	if err != nil {
		return nil, err
	}

	url = fmt.Sprintf("https://auth.levtech.jp%s", pathWithQuery)
	cookies, err = s.req3(client, url, cookies)
	if err != nil {
		return nil, err
	}

	pathParts := strings.Split(pathWithQuery, "/")
	pathID := pathParts[1]
	url = fmt.Sprintf("https://auth.levtech.jp/api/%s/signin-by-email", pathID)
	url, loginRes, cookies, err := s.req4(client, url, id, cookies)
	if err != nil {
		return nil, err
	}

	url, cookies, err = s.req5(client, loginRes.RedirectUri, cookies)
	if err != nil {
		return nil, err
	}

	url = "https://platform.levtech.jp/p/"
	url, cookies, err = s.req6(client, url, cookies)
	if err != nil {
		return nil, err
	}

	var session model.PlatformSession
	for _, cookie := range cookies {
		switch cookie.Name {
		case "CAKEPHP":
			session.SessionID = cookie.Value
		case "AWSELBAuthSessionCookie-0":
			session.AWSAuth = cookie.Value
		}
	}

	return &session, nil
}

// https://platform.levtech.jp/p/
func (s loginPlatformService) req1(client *http.Client, url string) (
	nextURL string, resCookies []*http.Cookie, err error,
) {
	fmt.Println("=== Request 1 ===")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Printf("Status Code: %d \n", resp.StatusCode)

	for k, v := range resp.Header {
		if strings.ToLower(k) == "location" && len(v) > 0 {
			nextURL = v[0]
		}
		fmt.Printf("%s: %s\n", k, v)
	}

	if resp.StatusCode >= 400 {
		err = fmt.Errorf(string(body))
		return
	}

	if resp.StatusCode < 300 {
		err = fmt.Errorf(
			"status code expected to be 3**, actual: %d", resp.StatusCode,
		)
		return
	}

	resCookies = resp.Cookies()
	err = nil
	return
}

// https://auth.levtech.jp/oidc/auth
func (s loginPlatformService) req2(
	client *http.Client, url string, cookies []*http.Cookie,
) (
	nextURL string, resCookies []*http.Cookie, err error,
) {
	fmt.Println("=== Request 2 ===")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Printf("Status Code: %d \n", resp.StatusCode)

	for k, v := range resp.Header {
		if strings.ToLower(k) == "location" && len(v) > 0 {
			nextURL = v[0]
		}
		fmt.Printf("%s: %s\n", k, v)
	}

	if resp.StatusCode >= 400 {
		err = fmt.Errorf(string(body))
		return
	}

	if resp.StatusCode < 300 {
		err = fmt.Errorf(
			"status code expected to be 3**, actual: %d", resp.StatusCode,
		)
		return
	}

	resCookies = append(cookies, resp.Cookies()...)
	err = nil
	return
}

// https://auth.levtech.jp/xxxx/signin?client_id=ltp
func (s loginPlatformService) req3(
	client *http.Client, url string, cookies []*http.Cookie,
) (
	resCookies []*http.Cookie, err error,
) {
	fmt.Println("=== Request 3 ===")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Printf("Status Code: %d \n", resp.StatusCode)

	if resp.StatusCode >= 400 {
		err = fmt.Errorf(string(body))
		return
	}

	if resp.StatusCode >= 300 {
		err = fmt.Errorf(
			"status code expected to be 2**, actual: %d", resp.StatusCode,
		)
		return
	}

	resCookies = append(cookies, resp.Cookies()...)
	err = nil
	return
}

type LoginResponse struct {
	Result      string `json:"result"`
	RedirectUri string `json:"redirectUri"`
}

// https://auth.levtech.jp/api/xxxx/signin-by-email
func (s loginPlatformService) req4(
	client *http.Client, url string,
	id model.PlatformID, cookies []*http.Cookie,
) (
	nextURL string, loginRes LoginResponse,
	resCookies []*http.Cookie, err error,
) {
	fmt.Println("=== Request 4 ===")

	var j []byte
	j, err = json.Marshal(id.RequestData())
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(j))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Printf("Status Code: %d \n", resp.StatusCode)

	for k, v := range resp.Header {
		if strings.ToLower(k) == "location" && len(v) > 0 {
			nextURL = v[0]
		}
		fmt.Printf("%s: %s\n", k, v)
	}

	if resp.StatusCode >= 400 {
		err = fmt.Errorf(string(body))
		return
	}

	err = json.Unmarshal(body, &loginRes)
	if err != nil {
		return
	}

	resCookies = append(cookies, resp.Cookies()...)
	return
}

// https://platform.levtech.jp//oauth2/idpresponse?code=xxx
func (s loginPlatformService) req5(
	client *http.Client, url string, cookies []*http.Cookie,
) (
	nextURL string, resCookies []*http.Cookie, err error,
) {
	fmt.Println("=== Request 5 ===")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Printf("Status Code: %d \n", resp.StatusCode)

	for k, v := range resp.Header {
		if strings.ToLower(k) == "location" && len(v) > 0 {
			nextURL = v[0]
		}
		fmt.Printf("%s: %s\n", k, v)
	}

	if resp.StatusCode >= 400 {
		err = fmt.Errorf(string(body))
		return
	}

	resCookies = append(cookies, resp.Cookies()...)
	err = nil
	return
}

// https://platform.levtech.jp/p/
func (s loginPlatformService) req6(
	client *http.Client, url string, cookies []*http.Cookie,
) (
	nextURL string, resCookies []*http.Cookie, err error,
) {
	fmt.Println("=== Request 6 ===")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Printf("Status Code: %d \n", resp.StatusCode)

	for k, v := range resp.Header {
		if strings.ToLower(k) == "location" && len(v) > 0 {
			nextURL = v[0]
		}
		fmt.Printf("%s: %s\n", k, v)
	}

	if resp.StatusCode >= 400 {
		err = fmt.Errorf(string(body))
		return
	}

	resCookies = append(cookies, resp.Cookies()...)
	err = nil
	return
}
