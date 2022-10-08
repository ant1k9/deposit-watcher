package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ant1k9/deposit-watcher/internal/datastruct"
)

const (
	// BaseURL is the only available URL `sravni.ru`
	BaseURL = "public.sravni.ru"
	// DepositsURI link for list of deposits. It has pagination
	DepositsURI = "/v1/deposit/special/list"
	// PageSize is a number of items for pagination
	PageSize = 10
)

func makeRequest(page int) (*http.Response, error) {
	client := http.Client{Timeout: 10 * time.Second}

	queryValues := url.Values{}
	queryValues.Add("location", "6.83")
	queryValues.Add("currency", "RUB")
	queryValues.Add("limit", strconv.Itoa(PageSize))
	queryValues.Add("skip", strconv.Itoa(PageSize*(page-1)))

	uri := url.URL{
		Scheme:   "https",
		Host:     BaseURL,
		Path:     DepositsURI,
		RawQuery: queryValues.Encode(),
	}

	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, err
	}

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
