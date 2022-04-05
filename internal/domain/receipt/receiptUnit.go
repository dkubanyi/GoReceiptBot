package receipt

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ReceiptUnit struct {
	CashRegisterCode           string `json:"cashRegisterCode"`
	BuildingNumber             int32  `json:"buildingNumber"`
	Country                    string `json:"country"`
	Municipality               string `json:"municipality"`
	PostalCode                 string `json:"postalCode"`
	PropertyRegistrationNumber string `json:"propertyRegistrationNumber"`
	StreetName                 string `json:"streetName"`
	Name                       string `json:"name"`
	UnitType                   string `json:"unitType"`
}

func (u ReceiptUnit) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func (u *ReceiptUnit) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &u)
}
