package receipt

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ReceiptOrganization struct {
	BuildingNumber             int32  `json:"buildingNumber"`
	Country                    string `json:"country"`
	Dic                        string `json:"dic"`
	IcDph                      string `json:"icDph"`
	Ico                        string `json:"ico"`
	Municipality               string `json:"municipality"`
	Name                       string `json:"name"`
	PostalCode                 string `json:"postalCode"`
	PropertyRegistrationNumber string `json:"propertyRegistrationNumber"`
	StreetName                 string `json:"streetName"`
	VatPayer                   bool   `json:"vatPayer"`
}

func (o ReceiptOrganization) Value() (driver.Value, error) {
	return json.Marshal(o)
}

func (o *ReceiptOrganization) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &o)
}
