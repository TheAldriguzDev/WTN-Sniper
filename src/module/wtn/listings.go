package wtn

import (
	error_manager "WTN-Sniper/src/packages/errors"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

type Listing struct {
	Pagination Pagination `json:"pagination"`
	Results    []Result   `json:"results"`
}

type Pagination struct {
	TotalPages   int64 `json:"totalPages"`
	Page         int64 `json:"page"`
	ItemsPerPage int64 `json:"itemsPerPage"`
	TotalItems   int64 `json:"totalItems"`
}

type Result struct {
	Product         Product `json:"product"`
	Name            string  `json:"name"`
	Price           int64   `json:"price"`
	IsLowestPrice   bool    `json:"isLowestPrice"`
	PaymentInfoUUID string  `json:"paymentInfoUuid"`
	PaymentOption   string  `json:"paymentOption"`
	VatSystem       string  `json:"vatSystem"`
	IsExpiringSoon  bool    `json:"isExpiringSoon"`
	ExpirationDate  string  `json:"expirationDate"`
	AddressUUID     string  `json:"addressUuid"`
	Image           string  `json:"image"`
}

type Product struct {
	Brand        string `json:"brand"`
	Image        string `json:"image"`
	Name         string `json:"name"`
	VariantID    int64  `json:"variantId"`
	EuropeanSize string `json:"europeanSize"`
	Sku          string `json:"sku"`
}

type Listings struct {
	Name     string
	Sku      string
	ID       string
	Size     string
	MinPrice string
}

func (s *WTNSession) exportListings() error {
	var tempList Listing
	rand.Seed(time.Now().UnixNano())

	limits := []string{"50", "100"}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api-sell.wethenew.com/listings?take=%s", limits[rand.Intn(len(limits))]), nil)
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

	if res.StatusCode != 200 {
		tempErr := error_manager.BadStatus(res.StatusCode, "Getting listings")
		return tempErr.Error
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tempErr := error_manager.ReadBody(err)
		return tempErr.Error
	}

	err = json.Unmarshal(body, &tempList)
	if err != nil {
		tempErr := error_manager.UnMarshal(err)
		return tempErr.Error
	}

	finalList := []Listings{}
	for _, product := range tempList.Results {
		finalList = append(finalList, Listings{
			Name: product.Product.Name,
			Sku:  product.Product.Sku,
			Size: product.Product.EuropeanSize,
			ID:   product.Name,
		})
	}

	err = saveListings(finalList)
	if err != nil {
		return fmt.Errorf("error while creating the csv. Open a ticket")
	}

	if !s.listingWebhook() {
		l.Error("Failed to send listing webhook!")
	}
	return nil
}

func saveListings(list []Listings) error {
	file, err := os.Create("./Settings/listings.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	heading := []string{"Name", "Sku", "Size", "MinPrice", "ID"}
	if err := writer.Write(heading); err != nil {
		return err
	}

	for _, product := range list {
		row := []string{product.Name, product.Sku, product.Size, "", product.ID}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func loadListings() ([]Listings, error) {
	file, err := os.Open("./Settings/listings.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	headings, err := reader.Read()
	if err != nil {
		return nil, err
	}

	fields := make(map[string]int)
	for i, campo := range headings {
		fields[campo] = i
	}

	var tempListings []Listings
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		tempData := Listings{
			Name:     row[fields["Name"]],
			Sku:      row[fields["Sku"]],
			Size:     row[fields["Size"]],
			MinPrice: row[fields["MinPrice"]],
			ID:       row[fields["ID"]],
		}

		tempListings = append(tempListings, tempData)
	}

	return tempListings, nil
}

func (s *WTNSession) extendListings() error {
	var tempList Listing
	rand.Seed(time.Now().UnixNano())

	limits := []string{"50", "100"}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api-sell.wethenew.com/listings?take=%s", limits[rand.Intn(len(limits))]), nil)
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

	if res.StatusCode != 200 {
		tempErr := error_manager.BadStatus(res.StatusCode, "Getting listings")
		return tempErr.Error
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tempErr := error_manager.ReadBody(err)
		return tempErr.Error
	}

	err = json.Unmarshal(body, &tempList)
	if err != nil {
		tempErr := error_manager.UnMarshal(err)
		return tempErr.Error
	}

	for _, product := range tempList.Results {
		err = s.extend(product)
		if err != nil {
			l.Error(err.Error())
		} else {
			l.Success(fmt.Sprintf("Time extended for %s", product.Product.Name))
			if s.extendWebhook(product) {
				l.Success("Log sent!")
			} else {
				l.Error("An error occured while loggin the webhook!")
			}
		}
	}
	return nil
}

type ExtendTime struct {
	ApplyToSameVariants bool `json:"applyToSameVariants"`
	Lifespan            int  `json:"lifespan"`
}

func (s *WTNSession) extend(product Result) error {
	ext := ExtendTime{
		ApplyToSameVariants: false,
		Lifespan:            60,
	}

	jsonPayload, err := json.Marshal(ext)
	if err != nil {
		tempErr := error_manager.Marshal(err)
		return tempErr.Error
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("https://api-sell.wethenew.com/listings/%s", product.Name), bytes.NewBuffer(jsonPayload))
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

	if res.StatusCode != 200 {
		tempErr := error_manager.BadStatus(res.StatusCode, "Extending time")
		return tempErr.Error
	}

	return nil
}
