package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"work_report/model"
)

type LoginPlatformService interface {
	Login(id model.PlatformID) (*model.PlatformSession, error)
}

type loginPlatformService struct {
	ParseHTMLService
}

func NewLoginPlatformService(ps ParseHTMLService) LoginPlatformService {
	return &loginPlatformService{ps}
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

	url, cookies, err = s.req2(client, url, cookies)
	if err != nil {
		return nil, err
	}

	cookies, scriptPath, err := s.req3(client, url, cookies)
	if err != nil {
		return nil, err
	}

	nextUUIDs, err := s.req4(client, scriptPath, cookies)
	if err != nil {
		return nil, err
	}

	url, loginRes, cookies, err := s.req5(client, url, id, nextUUIDs, cookies)
	if err != nil {
		return nil, err
	}

	url, cookies, err = s.req6(client, loginRes.RedirectUri, cookies)
	if err != nil {
		return nil, err
	}

	url = "https://platform.levtech.jp/p/"
	url, cookies, err = s.req7(client, url, cookies)
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
		err = fmt.Errorf("%s", string(body))
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
		err = fmt.Errorf("%s", string(body))
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
	resCookies []*http.Cookie, scriptPath string, err error,
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
		err = fmt.Errorf("%s", string(body))
		return
	}

	if resp.StatusCode >= 300 {
		err = fmt.Errorf(
			"status code expected to be 2**, actual: %d", resp.StatusCode,
		)
		return
	}

	ps := s.ParseHTMLService
	node, err := ps.Parse(string(body))
	if err != nil {
		return
	}

	node = ps.FindFirst(node, func(n *html.Node) bool {
		return ps.IsTag(n, "script") &&
			ps.AttrValRegExp(
				n, "src",
				"/_next/static/chunks/app/%5Buuid%5D/signin/page-[a-f0-9]+.js",
			)
	})
	scriptPath = ps.GetAttrValue(node, "src")

	resCookies = append(cookies, resp.Cookies()...)
	err = nil
	return
}

// https://auth.levtech.jp//_next/static/chunks/app/%5Buuid%5D/signin/page-xxx.js
func (s loginPlatformService) req4(
	client *http.Client, scriptPath string, cookies []*http.Cookie,
) (nextUUIDs []string, err error) {
	fmt.Println("=== Request 4 ===")

	url := fmt.Sprintf("https://auth.levtech.jp%s", scriptPath)

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
		err = fmt.Errorf("%s", string(body))
		return
	}

	if resp.StatusCode >= 300 {
		err = fmt.Errorf(
			"status code expected to be 2**, actual: %d", resp.StatusCode,
		)
		return
	}

	pattern := `[a-f0-9]{42}`
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}

	matched := re.FindAll(body, -1)
	if matched == nil {
		err = fmt.Errorf("uuid not found")
		return
	}

	nextUUIDs = make([]string, len(matched))
	for i, v := range matched {
		nextUUIDs[i] = string(v)
	}

	err = nil
	return
}

type LoginResponse struct {
	Result      string `json:"result"`
	RedirectUri string `json:"redirectUri"`
}

// https://auth.levtech.jp/xxxx/signin?client_id=ltp
func (s loginPlatformService) req5(
	client *http.Client, url string,
	id model.PlatformID, nextUUIDs []string, cookies []*http.Cookie,
) (
	nextURL string, loginRes LoginResponse,
	resCookies []*http.Cookie, err error,
) {
	fmt.Println("=== Request 5-1 ===")

	req1, err := http.NewRequest(
		http.MethodPost, url, bytes.NewBuffer([]byte(id.RequestData())),
	)
	if err != nil {
		return
	}

	req1.Header.Add("Content-Type", "text/plain")
	req1.Header.Add("Accept", "text/x-component")
	req1.Header.Add("Next-Action", nextUUIDs[0])

	for _, cookie := range cookies {
		req1.AddCookie(cookie)
	}

	resp, err := client.Do(req1)
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
		err = fmt.Errorf("%s", string(body))
		return
	}

	fmt.Println("=== Request 5-2 ===")

	req2, err := http.NewRequest(
		http.MethodPost, url, bytes.NewBuffer([]byte("[]")),
	)
	if err != nil {
		return
	}

	req2.Header.Add("Content-Type", "text/plain")
	req2.Header.Add("Accept", "text/x-component")
	req2.Header.Add("Next-Action", nextUUIDs[1])
	for _, cookie := range cookies {
		req2.AddCookie(cookie)
	}

	resp, err = client.Do(req2)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// 1:{"result":"SignedInCompletely","redirectUri":"https://platform.levtech.jp/oauth2/idpresponse?code=xxx&state=..."}
	// という行を探す
	found := false
	for _, line := range strings.Split(string(body), "\n") {
		if err := json.Unmarshal([]byte(line[2:]), &loginRes); err == nil {
			if loginRes.RedirectUri != "" {
				found = true
				break
			}
		}
	}

	if !found {
		err = fmt.Errorf("redirect uri not found")
		return
	}

	resCookies = append(cookies, resp.Cookies()...)
	return
}

// https://platform.levtech.jp//oauth2/idpresponse?code=xxx
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
		err = fmt.Errorf("%s", string(body))
		return
	}

	resCookies = append(cookies, resp.Cookies()...)
	err = nil
	return
}

// https://platform.levtech.jp/p/
func (s loginPlatformService) req7(
	client *http.Client, url string, cookies []*http.Cookie,
) (
	nextURL string, resCookies []*http.Cookie, err error,
) {
	fmt.Println("=== Request 7 ===")

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
		err = fmt.Errorf("%s", string(body))
		return
	}

	resCookies = append(cookies, resp.Cookies()...)
	err = nil
	return
}
