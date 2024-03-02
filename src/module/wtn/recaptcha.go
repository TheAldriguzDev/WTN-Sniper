package wtn

import (
	error_manager "WTN-Sniper/src/packages/errors"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

var siteKey string = "6LfbSlUpAAAAABNgkya850A9AtuIxEzJtv5V5cO5"

func (s *WTNSession) SolveCaptcha() (string, error) {
	captchaUrl := fmt.Sprintf("https://www.google.com/recaptcha/api2/anchor?ar=1&k=%s&co=aHR0cHM6Ly9zZWxsLndldGhlbmV3LmNvbTo0NDM.&hl=it&v=3sU2vDRVDmUU2E0Ro4VadvPr&size=invisible&cb=5up6oaiqzsc", siteKey)
	req, err := http.NewRequest(http.MethodGet, captchaUrl, nil)
	if err != nil {
		tempErr := error_manager.NewRequest(err)
		return "", tempErr.Error
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
		tempErr := error_manager.BadStatus(res.StatusCode, "Getting captcha challenge")
		return "", tempErr.Error
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tempErr := error_manager.ReadBody(err)
		return "", tempErr.Error
	}

	token := findToken(string(body))
	captchaResponse, err := s.bypassCaptcha(token)
	if err != nil {
		l.Error("Captcha token not found!")
		return "", err
	}

	return captchaResponse, nil
}

func findToken(responseText string) string {
	startDelimiter := `id="recaptcha-token" value="`
	endDelimiter := `">`

	startIndex := strings.Index(responseText, startDelimiter)
	if startIndex == -1 {
		return ""
	}

	endIndex := strings.Index(responseText[startIndex+len(startDelimiter):], endDelimiter)
	if endIndex == -1 {
		return ""
	}

	token := responseText[startIndex+len(startDelimiter) : startIndex+len(startDelimiter)+endIndex]

	return token
}

func (s *WTNSession) bypassCaptcha(token string) (string, error) {
	data := "v=UFwvoDBMjc8LiYc1DKXiAomK&reason=q&c=" + token + "&k=" + siteKey + "&co=aHR0cHM6Ly9hY2NvdW50cy5zcG90aWZ5LmNvbTo0NDM.&hl=en&size=invisible&chr=%5B61%2C36%2C84%5D&vh=7349404152&bg=!d3GgcVjIAAX6VNIG-kc72sZkL7AELV23BHEg3iiH1gcAAADHVwAAAAttAQecCimhpOYJBHsHw4TnDQnJAUU1KJWxkMVvr9kGAhPbfpEnRsIzZxDoK8WNA4Xk_jX6YLNl5cj97gy8xe0qj2UogYjr5xxWaD7OHCEWXqDqFHo9zQkvm1Jr-3PhDQbPfdz_WeOLnRGfdAlF7f6kTVJj8r_mdAx3g-11hZ4fXQpAMZ0qWUVIHOx4N86v_InW_G-9vhB6bzY_Xg1rQvsjsor6h9-BUi6cUZMmvYAn78v7JLBPZSdpWYD285rwy35stcDw5cYF8ruxzI_IqsNA6NAZWA4k1n-PuM5pxQDzLsrkD5oXB839hlcFKldmlsFx074KtmlcmvUVrD2O4Q8hNqjNTjRDSRqfZzcChfvqRKsx1DyhuXnz5dYAAR1ASd47CXlwOBdU6gCke1gRtLFtfSBMDvrhg7jK3uVK3jM-0q66IZwyUZHosVS0tI6DRRdXK6owLFJZi3lLnBzbASXdOaUGQHrzDFjQAbU76NE-neuga0bExaNraPqN0wYUBz1D0IPJ5kYLPQNArW2Z9a-to_yP1Oo7IDJlty6h9jTS7D32mQinK7JskejX_kPrchsfCCmNTQVmzSVnky-3WK4okaHXDKe9EHTTD2q5yQMRKNhHiHhifVwZ7fbBuaAP36m2qOxSrjrz6Em_HkCqAQb5GArJCVq7w04FuDxW7AGg086NIlp5QkST-AkJjCxn9BnV_5_37x_K8-vkMgIgND-pOet9HMf-4yrI7QrI8odYI9mmdEHwlSbKyWBkfnWFhTTL096dlgvDsNIiPoYIqdSjRAc1hkb2ToGqHTKD9VsfpcH79bBg_reEC7EK2ifubkSRhmz_LGhlN5wRTr_DuhjO0_pH8-TGKDCLlQJQ-jWC97z979drLh97I6wuXroq3xwfymQ9iDs2glksEExM78hbfVQsdRhLiUtDFkYxjinFEb325zUQJR6xT-yXNcLfLwDLWwfL6nuLS14IXKmENSW-6OIkXXJyDhgUKW6B1Reyll5b1s9A7OpsAtQH6H5rQzm6422zF12dO9JODF4UErQv43JpQu_wYG-VRUwGWcHvcrG9vp1c8mXjNfxoE2Ok0tTNXQXLr1DacKk4mG2YZF7X8xkPjDqW1XH6w8kde64MoCbMlc2u4yv9x-44P2XXMCDppFLsRLekxW00kAXP__rWpvNoEtt4PbI_Y_d9-lSLd6WDQ5mZObuIdo6BAS095B443-CNTb_4IAp9-4puXY_WU1cbkvt0hsV9iFjkcbsjAw1xmwBZoVg1ukewp4kPWL-oVVlGJYuQm_7AvAjZ6nRIRv_f7KebJQr-bY6wD3asqUzEZ8DHOLUJeScIFtDTFzAg9SxkP1dYde7y9umqn3a_3OyFR6iulqy-c0LoULRNh1DXG4KXKaabC4f7cixdSWPazY58wiic3ysahAsbaFGv_LzwFCy7uP0M7zKiwadGSOH_gaROuLTHbRnbEvPAgaa3zyP9mFPNhy_AsgKOAy5iA4l9qaBiVXwrWpXVyuQgsliVcmpeLSrMg9fpbb7LGcLv9dz5LxUetPIDUndRuJnW6xCyNakiVQMy6vF9l9qkEoHRjua7sPWZnJC29zjHdTgEVy5SOKinYGOBs1GlCFSSyIjBWixBWXH83hCjdd3TDJjQsNtDsRMr8mVlyiEqKkIttz1-2mV2ZhA6FmJ3Ldm-tnlQN5iaIM40HKbbrHDuDKWhdWXsouO2BfLJxDDvu--e251eTYOuFIHRQzCUg-y4LddffPqdFpemjCsJC1xHTx2DQXZTdQTv9n0FK9GTwRizJxnYqn5lDoXZv4MtG4tSZFjx8U9KnNBApFcDXXfFeymVWT5miPlimr9zSBRGEdzCAv-NctpVXrSwd3Yzpsj_eGFT91owgeAzjOnPFWMod2XCZbEywubJI-0QsHFxwiGbsnXtV2fXQOzdpdQKVcynj7gQhIJHQogB0M4achR6TmT-7dRvNkffa17qyRqpoIXhbpbvC9cgQG4VQaYjlhpeiWNHM7uTW1-5cdkZOcVqxsU5c1fzMpv77BuRYm--EJpFZSfsihySlvcGWVnz1qS-deD4gGUa8un8j3v0-YAu4llS6vC9OpCb-khnh7SgHk-a19cLD9m6mXVu9EJlV2gbdMcKouobkIljeKBT6ivhkemTe_peKfgDjgFSJfJ7Hxey2LR7nG1YW2FVv6kAOPRfoHNUf2OEvUHcvZ00jc7nJZTfrt-8nmutFD9C59MQ5HWvtIK5XobaAxyunaZon6iZxiFFRDo2o-xl6TwCuHYvmVWl6mAr5kn5QDlclIKc6hrIq8osYCcukWMZhu7L9wsyVMy1WC2GhXdWlTZnaJjqLtGBsxaTCbzND4nZ0zGEsGMX180J-y1PQ3EY3nP0e4ToqO8rXPi6lZ4GmGTpm0XypZ0jkf1xnU1FacQhmpVmIKru8kbjjChfywMM2exkn3E7CINxQS77i81vn3c8fWcdvKQ9lVProo60Yzea7RpjOdnfk9T4CcjV-J941093qAttWyknhB661xBQCzOXFB0euSb8Jn-J_5tSgX4NE1AyNXQEA5wk6km6tT3UUyK3yTEn8oynK_FZz_p4W4BGy_sCUm0IG43ioT-17L2CoQAzk2ZE5g4eh7jkASVHBeXREbMWtB4YdO-gPwxIrWVVOiN57jSDi5yM08wgBqKeAVYLHXFFuUG7konyayI7tTwxYjN0j7T9nGR2Jh1wmA-q99D4tbsM2AvWIWn9j3g83JBF4nqzS4lt72WUpL3kAdbOz2xwRKaWLFaEsaM9jQeg2ijJpTNqRlKxtXneWqjkca5JZCZEmGCbplWJAEARNOEVHWd00dc2dCt8KEHiBiAP86L6loq_QvD-kbLd1bd9S8FqCMVFRcOwOOBvUEBm1D-mJiA2KWBJ9T87kcAQmLRQxrTuGHMojr9cBtKz-2afsMXRPoCPmRc-dDwiYOXUgdERgEH6lifStYMTZcjS66GGA-0UccFdY2yAl7TG6b3lDa-lbTwSJHESj_UrH3neTMf2U8Z6rFWsTIHa5XfQ8nFacgUyokFLtzxGH57QQqRUc0bfqEouu_o8S0galOM1p2uaZqrrdvAbq31i-xMU5CqW0_WVG3REfA6SY0CJLXOs8mzwGFgZJEpr374MMRL6JEUu7qd_jib4P9-O8pvKFk7tfPTccXWq12b1gj7SsA6sdeffMMG1gpD-kYGud8ghD6x9sevkZ-IRveRZQmUCqXvT6rl-YOfyBTDsv2vpqD1kXxGSNV206XBFw6bFQB583TBhFWfm3p6nc1s3p-KY4oIMR1l6Z5Ccfh7CWv7EYNkbjwfsrk1PXoI38vy4cT8ttz49TQ5WSPSBgeZuAKUlX0Hml2C2xtis_a3YABvB4UsJK65Rg7hQCWLAlX8HYLeVYiUiqh31LE5JUPayYiC0nxQADw7A6-6hFtJqQz84LrxMw7a-Q59R0VwCqWGfCebmRh_BVlg"
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://www.google.com/recaptcha/api2/reload?k=%s", siteKey), strings.NewReader(data))
	if err != nil {
		tempErr := error_manager.NewRequest(err)
		return "", tempErr.Error
	}

	req.Header = http.Header{
		"content-type": {"application/x-www-form-urlencoded"},
	}

	res, err := s.HttpClient.TLS_client.Do(req)
	if err != nil {
		tempErr := error_manager.DoRequest(err)
		return "", tempErr.Error
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tempErr := error_manager.ReadBody(err)
		return "", tempErr.Error
	}

	newBody := string(body)
	index := strings.Index(newBody, "[")
	if index != -1 {
		finalBody := newBody[index:]
		var response []interface{}
		err = json.Unmarshal([]byte(finalBody), &response)
		if err != nil {
			tempErr := error_manager.UnMarshal(err)
			return "", tempErr.Error
		}

		return response[1].(string), nil
	} else {
		return "", fmt.Errorf("captcha response not found")
	}
}
