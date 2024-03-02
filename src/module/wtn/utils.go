package wtn

import (
	"fmt"

	http "github.com/bogdanfinn/fhttp"
)

func (s *WTNSession) getHeaders() http.Header {
	headers := http.Header{
		"accept":             {"application/json, text/plain, */*"},
		"accept-language":    {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"authority":          {"api-sell.wethenew.com"},
		"authorization":      {fmt.Sprintf("Bearer %s", s.loginData.User.AccessToken)},
		"origin":             {"https://sell.wethenew.com"},
		"pragma":             {"no-cache"},
		"referer":            {"https://sell.wethenew.com/"},
		"sec-ch-ua":          {"\"Google Chrome\";v=\"119\"}, \"Chromium\";v=\"119\"}, \"Not?A_Brand\";v=\"24\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"macOS\""},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-site"},
		"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"},
		"x-xss-protection":   {"1;mode=block"},
	}

	return headers
}
