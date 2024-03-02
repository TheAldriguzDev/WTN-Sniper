package wtn

import (
	error_manager "WTN-Sniper/src/packages/errors"
	"encoding/json"
	"io"

	http "github.com/bogdanfinn/fhttp"
)

type Csrf struct {
	CSRF_Token string `json:"csrfToken"`
}

func (s *WTNSession) getCsrf() (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://sell.wethenew.com/api/auth/csrf", nil)
	if err != nil {
		tempErr := error_manager.NewRequest(err)
		return "", tempErr.Error
	}

	req.Header = http.Header{
		"authority":          {"sell.wethenew.com"},
		"accept":             {"*/*"},
		"accept-language":    {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"content-type":       {"application/json"},
		"if-none-match":      {"W/\"v9kmf4kz2d28\""},
		"referer":            {"https://sell.wethenew.com/login"},
		"sec-ch-ua":          {"\"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"macOS\""},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-origin"},
		"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"},
	}

	res, err := s.HttpClient.TLS_client.Do(req)
	if err != nil {
		tempErr := error_manager.DoRequest(err)
		return "", tempErr.Error
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		tempErr := error_manager.BadStatus(res.StatusCode, "CSRF token")
		return "", tempErr.Error
	}

	var tempBody Csrf
	body, err := io.ReadAll(res.Body)
	if err != nil {
		tempErr := error_manager.ReadBody(err)
		return "", tempErr.Error
	}

	err = json.Unmarshal(body, &tempBody)
	if err != nil {
		tempErr := error_manager.UnMarshal(err)
		return "", tempErr.Error
	}

	return tempBody.CSRF_Token, nil
}
