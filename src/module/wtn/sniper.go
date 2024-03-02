package wtn

import (
	error_manager "WTN-Sniper/src/packages/errors"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

func (s *WTNSession) sniper() error {
	listings, err := loadListings()
	if err != nil {
		l.Error(fmt.Sprintf("Error while loading the listings! %s", err.Error()))
		time.Sleep(20 * time.Second)
		os.Exit(0)
	}

	if len(listings) <= 0 {
		return fmt.Errorf("Product list require to run the sniper")
	}

	var tempOffer Offer
	var wg sync.WaitGroup

	wg.Add(1)
	go func(tempOffer *Offer) {
		defer wg.Done()

		for {
			err := s.getOffers(tempOffer)
			if err != nil {
				if err.Error() == "login session expired" {
					l.Warning(err.Error())
					// Create the new session
					// WTN module has been initialised
					loginStep := true
					for loginStep {
						loginSession, err := s.Login()
						if err != nil {
							l.Error(err.Error())
							time.Sleep(time.Duration(s.settingsData.Delay) * time.Second)
						}
						s.loginData = loginSession
						loginStep = false
						break
					}

					// Save login data
					err = s.saveLoginData()
					if err != nil {
						l.Error("An error occured while saving the login data! Skipping...")
					}
					l.Success("Logged in!")
				} else {
					l.Error(err.Error())
				}
			}
			time.Sleep(time.Duration(s.settingsData.Delay) * time.Second)
		}
	}(&tempOffer)

	wg.Add(1)
	go func(listings []Listings, tempOffer *Offer) {
		defer wg.Done()

		for {
			err := s.snipeOffer(listings, tempOffer)
			if err != nil {
				if err.Error() != "no offers found" {
					l.Error(err.Error())
				}
			}
			// Temp solution to avoid loosing offers
			if s.settingsData.Delay >= 1 {
				time.Sleep(500 * time.Millisecond)
			} else {
				time.Sleep(350 * time.Millisecond)
			}
		}
	}(listings, &tempOffer)

	wg.Wait()

	return nil
}

/*
- Go routine that monitors the offers
- Go routine that confirm the offers

- Communication with a struct pointed from both the functions
*/

type Offer struct {
	Pagination Offer_Pagination `json:"pagination"`
	Results    []Offer_Result   `json:"results"`
}

type Offer_Pagination struct {
	TotalPages   int64 `json:"totalPages"`
	Page         int64 `json:"page"`
	ItemsPerPage int64 `json:"itemsPerPage"`
	TotalItems   int64 `json:"totalItems"`
}

type Offer_Result struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	VariantID    int64  `json:"variantId"`
	Sku          string `json:"sku"`
	Brand        string `json:"brand"`
	Image        string `json:"image"`
	EuropeanSize string `json:"europeanSize"`
	ListingPrice int64  `json:"listingPrice"`
	Price        int64  `json:"price"`
	CreateTime   string `json:"createTime"`
}

func (s *WTNSession) getOffers(tempOffer *Offer) error {
	l.General("Waiting for offers...")

	_, err := s.RotateProxy()
	if err != nil {
		l.Error(fmt.Sprintf("An error occured while rotating the proxy: %s", err.Error()))
	}

	rand.Seed(time.Now().UnixNano())
	limits := []string{"50", "100"}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api-sell.wethenew.com/offers?take=%s", limits[rand.Intn(len(limits))]), nil)
	if err != nil {
		tempErr := error_manager.NewRequest(err)
		return tempErr.Error
	}
	req.Header = s.getHeaders()

	res, err := s.HttpClient.TLS_client.Do(req)
	if err != nil {
		tempErr := error_manager.DoRequest(err)
		return tempErr.Error
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case 401:
		return fmt.Errorf("login session expired")
	case 200:
		break
	default:
		tempErr := error_manager.BadStatus(res.StatusCode, "Monitoring offers")
		return tempErr.Error
	}
	// If we are here the response was 200

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tempErr := error_manager.ReadBody(err)
		return tempErr.Error
	}

	err = json.Unmarshal(body, &tempOffer) // Doing so the "temppOffer" will be overwritten every time --> we could miss an offer due to this
	if err != nil {
		tempErr := error_manager.UnMarshal(err)
		return tempErr.Error
	}

	return nil
}

func (s *WTNSession) snipeOffer(listing []Listings, tempOffer *Offer) error {
	// First of all we need to parse the offer

	if len(tempOffer.Results) <= 0 {
		return fmt.Errorf("no offers found")
	}

	for _, offer := range tempOffer.Results {
		go s.offerFoundWebhook(offer) // Go routine to remove the delay of the discord request
		if checkOffer(listing, offer) {
			// Offer needs to be sniped
			err := s.confirmOffer(offer)
			if err != nil {
				l.Error(fmt.Sprintf("An error occured while sniping the offer: %e", err))
				if s.acceptErrorWebhook(offer, err.Error()) {
					l.Success("Webhook sent")
				} else {
					l.Error("An error occured while sending the webhook")
				}
			} else {
				l.Success(fmt.Sprintf("Offer %s accepted!", offer.Name))
				if s.acceptWebhook(offer) {
					l.Success("Webhook sent")
				} else {
					l.Error("An error occured while sending the webhook")
				}
			}
		} else {
			l.Error(fmt.Sprintf("Offer for %s is under min price!", offer.Name))
		}
	}

	return nil
}

type ConfirmOffer struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	VariantId int64  `json:"variantId"`
}

func (s *WTNSession) confirmOffer(offer Offer_Result) error {
	tempOrder := ConfirmOffer{
		Name:      offer.ID,
		Status:    "ACCEPTED",
		VariantId: offer.VariantID,
	}

	jsonPayload, err := json.Marshal(tempOrder)
	if err != nil {
		tempErr := error_manager.Marshal(err)
		return tempErr.Error
	}

	req, err := http.NewRequest(http.MethodPost, "https://api-sell.wethenew.com/offers", bytes.NewBuffer(jsonPayload))
	if err != nil {
		tempErr := error_manager.NewRequest(err)
		return tempErr.Error
	}

	req.Header = http.Header{
		"Accept":             {"application/json, text/plain, */*"},
		"Authorization":      {fmt.Sprintf("Bearer %s", s.loginData.User.AccessToken)},
		"Content-Type":       {"application/json"},
		"Pragma":             {"no-cache"},
		"Referer":            {"https://sell.wethenew.com/"},
		"User-Agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"},
		"X-XSS-Protection":   {"1;mode=block"},
		"sec-ch-ua":          {"\"Google Chrome\";v=\"119\"}, \"Chromium\";v=\"119\"}, \"Not?A_Brand\";v=\"24\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"macOS\""},
	}

	res, err := s.HttpClient.TLS_client.Do(req)
	if err != nil {
		tempErr := error_manager.DoRequest(err)
		return tempErr.Error
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case 200, 201, 204:
		return nil
	default:
		tempErr := error_manager.BadStatus(res.StatusCode, "Sniping the offer")
		return tempErr.Error
	}
}

// Return true if the offer needs to be accepted
// Return false if the offer is to ignore
func checkOffer(listing []Listings, offer Offer_Result) bool {
	for _, p := range listing {
		if p.ID == offer.ID {
			tempPrice, _ := strconv.ParseInt(p.MinPrice, 10, 64)
			if offer.Price > tempPrice {
				return true
			}
		}
	}

	return false
}
