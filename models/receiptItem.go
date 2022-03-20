package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ReceiptItem struct {
	Name     string  `json:"name"`
	ItemType string  `json:"itemType"`
	Quantity float64 `json:"quantity"`
	VatRate  float64 `json:"vatRate"`
	Price    float64 `json:"price"`
}

type ReceiptItems []ReceiptItem

func (s ReceiptItems) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *ReceiptItems) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	}
	return errors.New("type assertion failed")
}
