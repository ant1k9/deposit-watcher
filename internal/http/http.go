package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ant1k9/deposit-watcher/internal/datastruct"
)

const (
	// BaseURL is the only available URL `sravni.ru`
	BaseURL = "https://www.sravni.ru"
	// DepositsURI link for list of deposits. It has pagination
	DepositsURI = "/proxy-deposits/deposits/list"
	// PageSize is a number of items for pagination
	PageSize = 10
)

func makeRequest(page int) (*http.Response, error) {
	client := http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest(http.MethodPost, BaseURL+DepositsURI, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json")
	req.Body = ioutil.NopCloser(strings.NewReader(
		fmt.Sprintf(`{
			"filters": {
				"organization": [],
				"additionalConditions": [],
				"prolongation": [],
				"interestPaymentMethod": [],
				"capitalization": [],
				"currency": "RUB",
				"earlyTermination": [],
				"depositTypes": [],
				"interestPayment": [],
				"rating": "100",
				"location":"6.83.",
				"advertising": {
					"source": "search"
				}
			},
			"limit": %d, "skip":%d
		}`, PageSize, PageSize*(page-1)),
	))

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetDeposits receives all deposits on a page `page` with size PageSize
func GetDeposits(page int) []datastruct.Deposit {
	var single []datastruct.Deposit
	deposits := make([]datastruct.Deposit, 0)

	entries := struct {
		Data json.RawMessage `json:"items"`
	}{}

	response, err := makeRequest(page)

	if err == nil {
		data, err := ioutil.ReadAll(response.Body)
		if err == nil {
			_ = json.Unmarshal(data, &entries)
			_ = json.Unmarshal(entries.Data, &single)
			deposits = append(deposits, single...)

			for _, deposit := range single {
				var group []datastruct.Deposit
				_ = json.Unmarshal(deposit.Other, &group)
				deposits = append(deposits, group...)
			}
		}
	}

	return deposits
}

// GetTotal receives total number of banks on the site
func GetTotal(page int) int {
	var total datastruct.Total
	response, err := makeRequest(page)

	if err == nil {
		data, err := ioutil.ReadAll(response.Body)
		if err == nil {
			_ = json.Unmarshal(data, &total)
		}
	}

	return total.Total
}
