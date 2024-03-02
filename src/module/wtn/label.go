package wtn

import (
	error_manager "WTN-Sniper/src/packages/errors"
	"encoding/json"
	"fmt"
	"io"
	"os"

	http "github.com/bogdanfinn/fhttp"
)

type Label struct {
	Label_Pagination Label_Pagination `json:"Label_Pagination"`
	Label_Results    []Label_Result   `json:"Label_Results"`
}

type Label_Pagination struct {
	TotalPages   int64 `json:"totalPages"`
	Page         int64 `json:"page"`
	ItemsPerPage int64 `json:"itemsPerPage"`
	TotalItems   int64 `json:"totalItems"`
}

type Label_Result struct {
	CreateTime    string        `json:"createTime"`
	Name          string        `json:"name"`
	Status        string        `json:"status"`
	Label_Product Label_Product `json:"Label_Product"`
}

type Label_Product struct {
	Brand string `json:"brand"`
	Name  string `json:"name"`
	Sku   string `json:"sku"`
	Image string `json:"image"`
	Size  string `json:"size"`
	Price int64  `json:"price"`
}

func (s *WTNSession) label() error {
	req, err := http.NewRequest(http.MethodGet, "https://api-sell.wethenew.com/deals?skip=0&take=100&dealStatuses%5B%5D=TO_RECEIVE_SHIP", nil)
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
		tempErr := error_manager.BadStatus(res.StatusCode, "Getting sales")
		return tempErr.Error
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tempErr := error_manager.ReadBody(err)
		return tempErr.Error
	}

	var tempLabel Label
	err = json.Unmarshal(body, &tempLabel)
	if err != nil {
		tempErr := error_manager.UnMarshal(err)
		return tempErr.Error
	}

	// We can now parse and download every label

	if len(tempLabel.Label_Results) <= 0 {
		return fmt.Errorf("no label to download found")
	}

	for _, sale := range tempLabel.Label_Results {
		err = s.downloadLabel(sale)
		if err != nil {
			l.Error(err.Error())
		} else {
			l.Success(fmt.Sprintf("Label %s downloaded", sale.Name))
			if s.labelWebhook(sale) {
				l.Success("Webhook sent!")
			} else {
				l.Error("An error occured while loggin the webhook!")
			}
		}
	}
	return nil
}

func (s *WTNSession) downloadLabel(sale Label_Result) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api-sell.wethenew.com/deals/documents/%s", sale.Name), nil)
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
		tempErr := error_manager.BadStatus(res.StatusCode, "Download label")
		return tempErr.Error
	}

	file, err := os.Create(fmt.Sprintf("./Data/%s%s.pdf", sale.Label_Product.Name, sale.Label_Product.Size))
	if err != nil {
		return fmt.Errorf("error while creating the pdf for %s", sale.Label_Product.Name)
	}
	defer file.Close()
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return fmt.Errorf("error while downloading the pdf for %s", sale.Label_Product.Name)
	}

	return nil
}
