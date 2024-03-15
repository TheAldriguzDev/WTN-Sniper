package wtn

import (
	error_manager "WTN-Sniper/src/packages/errors"
	"WTN-Sniper/src/packages/session"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	http "github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/tls-client/profiles"
)

func (s *WTNSession) Login() (LoginSession, error) {
	var tempSession LoginSession
	l.General("Starting login...")

	// First of all we check if we already have a saved session
	loginStatus, err := s.loadLoginData()
	if err != nil {
		l.Error(err.Error())
	} else if loginStatus {
		l.Success("Session found!")
		tempSession = s.loginData
		return tempSession, nil
	}

	err = s.credentials()
	if err != nil {
		return tempSession, err
	}

	sess, err := s.session()
	if err != nil {
		return sess, err
	}

	s.saveLoginData()

	return sess, err
}

func (s *WTNSession) credentials() error {
	csrf, err := s.getCsrf()
	if err != nil {
		return err
	}

	l.General("Solving captcha challenge...")
	captchaResult, err := s.SolveCaptcha()
	if err != nil {
		os.Exit(0)
	}

	redirect := "redirect=false"
	email := "&email=" + url.QueryEscape(s.settingsData.WTN.Email)
	password := "&password=" + url.QueryEscape(s.settingsData.WTN.Password)
	recaptchaToken := "&recaptchaToken=" + url.QueryEscape(captchaResult)
	pushToken := "&pushToken=" + url.QueryEscape("undefined")
	os := "&os=" + url.QueryEscape("undefined")
	osVersion := "&osVersion=" + url.QueryEscape("undefined")
	csrfToken := "&csrfToken=" + url.QueryEscape(csrf)
	callbackUrl := "&callbackUrl=" + url.QueryEscape("https://sell.wethenew.com/login")
	json := "&json=true"
	payload := redirect + email + password + recaptchaToken + pushToken + os + osVersion + csrfToken + callbackUrl + json

	l.General("Submitting login...")
	req, err := http.NewRequest(http.MethodPost, "https://sell.wethenew.com/api/auth/callback/credentials", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		tempErr := error_manager.NewRequest(err)
		return tempErr.Error
	}

	req.Header = http.Header{
		"authority":          {"sell.wethenew.com"},
		"accept":             {"*/*"},
		"accept-language":    {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"content-type":       {"application/x-www-form-urlencoded"},
		"origin":             {"https://sell.wethenew.com"},
		"referer":            {"https://sell.wethenew.com/login"},
		"sec-ch-ua":          {"\"Google Chrome\";v=\"119\"}, \"Chromium\";v=\"119\"}, \"Not?A_Brand\";v=\"24\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"macOS\""},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-origin"},
		"traceparent":        {"00-50a3707a4dbbe34ef9154bc6d2a998d0-0ef8793b794765be-01"},
		"tracestate":         {"3403653@nr=0-1-3403653-322552406-0ef8793b794765be----1707326181526"},
		"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"},
	}

	res, err := s.HttpClient.TLS_client.Do(req)
	if err != nil {
		tempErr := error_manager.DoRequest(err)
		return tempErr.Error
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tempErr := error_manager.ReadBody(err)
		return tempErr.Error
	}

	if string(body) == `{"url":"https://sell.wethenew.com/login"}` {
		return nil
	} else {
		return fmt.Errorf("error while loggin in! %s", string(body))
	}
}

type LoginSession struct {
	User    User           `json:"user"`
	Expires string         `json:"expires"`
	Cookies []*http.Cookie `json:"cookies"`
}

type User struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (s *WTNSession) session() (LoginSession, error) {
	l.General("Getting session...")
	var tempSession LoginSession
	req, err := http.NewRequest(http.MethodGet, "https://sell.wethenew.com/api/auth/session", nil)
	if err != nil {
		tempErr := error_manager.NewRequest(err)
		return tempSession, tempErr.Error
	}

	req.Header = http.Header{
		"authority":          {"sell.wethenew.com"},
		"accept":             {"*/*"},
		"accept-language":    {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"content-type":       {"application/json"},
		"if-none-match":      {"\"bwc9mymkdm2\""},
		"referer":            {"https://sell.wethenew.com/login"},
		"sec-ch-ua":          {"\"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"macOS\""},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-origin"},
		"traceparent":        {"00-eee66496aec431055f70586863e385b0-6c459e74c4f24e52-01"},
		"tracestate":         {"3403653@nr=0-1-3403653-322552406-6c459e74c4f24e52----1707326182625"},
		"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"},
	}

	res, err := s.HttpClient.TLS_client.Do(req)
	if err != nil {
		tempErr := error_manager.DoRequest(err)
		return tempSession, tempErr.Error
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 && res.StatusCode != 304 {
		tempErr := error_manager.BadStatus(res.StatusCode, "Getting the session")
		return tempSession, tempErr.Error
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tempErr := error_manager.ReadBody(err)
		return tempSession, tempErr.Error
	}

	err = json.Unmarshal(body, &tempSession)
	if err != nil {
		tempErr := error_manager.UnMarshal(err)
		return tempSession, tempErr.Error
	}

	cookies := s.HttpClient.TLS_client.GetCookies(res.Request.URL)
	tempSession.Cookies = cookies

	return tempSession, nil
}

type loginData struct {
	LoginSession LoginSession `json:"login_session"`
}

func (s *WTNSession) saveLoginData() error {
	tempData := loginData{
		LoginSession: s.loginData,
	}

	jsonData, err := json.MarshalIndent(tempData, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./Data/login.json", jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *WTNSession) loadLoginData() (bool, error) {
	loginData := loginData{}

	content, err := ioutil.ReadFile("./Data/login.json")
	if err != nil {
		return false, nil
	}

	err = json.Unmarshal(content, &loginData)
	if err != nil {
		tempErr := error_manager.UnMarshal(err)
		return false, tempErr.Error
	}

	if (loginData.LoginSession.User.AccessToken == "") || (loginData.LoginSession.User.RefreshToken == "") {
		return false, nil
	}

	s.loginData = loginData.LoginSession

	_, err = s.refreshSession()
	if err != nil {
		// Return to create a new login session
		l.Error(err.Error())
		return false, nil
	}
	err = s.saveLoginData()
	if err != nil {
		l.Error(err.Error())
		return false, err
	}

	err = s.checkSession()
	if err != nil {
		l.Warning("Login session not found!")
		return false, nil
	}

	return true, nil
}

func (s *WTNSession) refreshSession() (LoginSession, error) {
	var tempSession LoginSession
	req, err := http.NewRequest(http.MethodGet, "https://sell.wethenew.com/api/auth/session", nil)
	if err != nil {
		tempErr := error_manager.NewRequest(err)
		return tempSession, tempErr.Error
	}

	var cookie_message string
	for _, cookie := range s.loginData.Cookies {
		cookie_message = cookie_message + fmt.Sprintf("%s=%s; ", cookie.Name, cookie.Value)
	}

	req.Header = http.Header{
		"authority":          {"sell.wethenew.com"},
		"accept":             {"*/*"},
		"accept-language":    {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"content-type":       {"application/json"},
		"cookie":             {cookie_message},
		"if-none-match":      {"\"bwc9mymkdm2\""},
		"newrelic":           {"eyJ2IjpbMCwxXSwiZCI6eyJ0eSI6IkJyb3dzZXIiLCJhYyI6IjM0MDM2NTMiLCJhcCI6IjMyMjU1MjQwNiIsImlkIjoiNmM0NTllNzRjNGYyNGU1MiIsInRyIjoiZWVlNjY0OTZhZWM0MzEwNTVmNzA1ODY4NjNlMzg1YjAiLCJ0aSI6MTcwNzMyNjE4MjYyNX19"},
		"referer":            {"https://sell.wethenew.com/login"},
		"sec-ch-ua":          {"\"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"macOS\""},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-origin"},
		"traceparent":        {"00-eee66496aec431055f70586863e385b0-6c459e74c4f24e52-01"},
		"tracestate":         {"3403653@nr=0-1-3403653-322552406-6c459e74c4f24e52----1707326182625"},
		"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"},
	}

	res, err := s.HttpClient.TLS_client.Do(req)
	if err != nil {
		tempErr := error_manager.DoRequest(err)
		return tempSession, tempErr.Error
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		tempErr := error_manager.BadStatus(res.StatusCode, "Refreshing the session")
		return tempSession, tempErr.Error
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tempErr := error_manager.ReadBody(err)
		return tempSession, tempErr.Error
	}

	err = json.Unmarshal(body, &tempSession)
	if err != nil {
		tempErr := error_manager.UnMarshal(err)
		return tempSession, tempErr.Error
	}

	cookies := s.HttpClient.TLS_client.GetCookies(res.Request.URL)
	tempSession.Cookies = cookies

	return tempSession, nil

}

func (s *WTNSession) checkSession() error {
	req, err := http.NewRequest(http.MethodGet, "https://api-sell.wethenew.com/sellers/me", nil)
	if err != nil {
		tempErr := error_manager.NewRequest(err)
		return tempErr.Error
	}

	req.Header = s.getHeaders()

	var cookie_message string
	for _, cookie := range s.loginData.Cookies {
		cookie_message = cookie_message + fmt.Sprintf("%s=%s; ", cookie.Name, cookie.Value)
	}

	req.Header.Add("cookie", cookie_message)

	res, err := s.HttpClient.TLS_client.Do(req)
	if err != nil {
		tempErr := error_manager.DoRequest(err)
		return tempErr.Error
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 && res.StatusCode != 304 {
		tempErr := error_manager.BadStatus(res.StatusCode, "Checking login session")
		s.loginData = LoginSession{}
		newSession, err := session.New(30, profiles.Chrome_120, s.Proxy_Manager)
		if err != nil {
			l.Error(err.Error())
			time.Sleep(20 * time.Second)
			os.Exit(1)
		}
		s.Session = newSession
		return tempErr.Error
	}

	return nil
}
