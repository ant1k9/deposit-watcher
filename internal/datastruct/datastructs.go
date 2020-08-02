package datastruct

import (
	"encoding/json"
)

// Deposit represent description with JSON tags
type Deposit struct {
	Alias  string `json:"alias"`
	Name   string `json:"name"`
	Amount struct {
		From int `json:"from"`
	} `json:"amount"`
	Rate         float64 `json:"rate"`
	Organization struct {
		Alias string `json:"alias"`
		Name  string `json:"name"`
	} `json:"organization"`
	Replenishment struct {
		Available   bool   `json:"available,omitempty"`
		Description string `json:"description,omitempty"`
	} `json:"replenishment"`
	Other json.RawMessage `json:"group,omitempty"`
}

// Total is a partial structure that only gets a total number
// of banks from the response.
type Total struct {
	Total int `json:"total"`
}

// BankRow is a struct to get bank records from database
type BankRow struct {
	ID    int    `db:"id" json:"id"`
	Alias string `db:"alias" json:"alias"`
	Name  string `db:"name" json:"name"`
}

// DepositRow is a struct to get deposit records from database
type DepositRow struct {
	ID               int     `db:"id" json:"id"`
	Alias            string  `db:"alias" json:"alias"`
	BankID           int     `db:"bank_id" json:"bankId"`
	Detail           string  `db:"detail" json:"detail"`
	HasReplenishment bool    `db:"has_replenishment" json:"hasReplenishment"`
	IsUpdated        bool    `db:"is_updated" json:"isUpdated"`
	MinimalAmount    int     `db:"minimal_amount" json:"minimalAmount"`
	Name             string  `db:"name" json:"name"`
	Off              bool    `db:"off" json:"off"`
	IsExist          bool    `db:"is_exist" json:"isExist"`
	PreviousRate     float64 `db:"previous_rate" json:"previousRate"`
	Rate             float64 `db:"rate" json:"rate"`
	UpdatedAt        string  `db:"updated_at" json:"updatedAt"`
}

type DepositRowShort struct {
	ID               int     `db:"id" json:"id"`
	Alias            string  `db:"alias" json:"alias"`
	BankName         string  `db:"bank_name" json:"bankName"`
	HasReplenishment bool    `db:"has_replenishment" json:"hasReplenishment"`
	Description      string  `db:"detail" json:"description"`
	Name             string  `db:"name" json:"name"`
	Rate             float64 `db:"rate" json:"rate"`
}
