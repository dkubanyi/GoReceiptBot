package entities

// FinancnaSpravaResponse /**
//* curl --location --request POST 'https://ekasa.financnasprava.sk/mdu/api/v1/opd/receipt/find' \
//--header 'Content-Type: application/json' \
//--data-raw '{
//    "receiptId": "O-863A291A95A64D7EBA291A95A63D7E4C"
//}'
type FinancnaSpravaResponse struct {
	ReturnValue int `json:"returnValue"`
	Receipt     struct {
		ReceiptID        string      `json:"receiptId"`
		Ico              string      `json:"ico"`
		CashRegisterCode string      `json:"cashRegisterCode"`
		IssueDate        string      `json:"issueDate"`
		CreateDate       string      `json:"createDate"`
		CustomerID       interface{} `json:"customerId"`
		Dic              string      `json:"dic"`
		IcDph            string      `json:"icDph"`
		InvoiceNumber    interface{} `json:"invoiceNumber"`
		Okp              string      `json:"okp"`
		Paragon          bool        `json:"paragon"`
		ParagonNumber    interface{} `json:"paragonNumber"`
		Pkp              string      `json:"pkp"`
		ReceiptNumber    int         `json:"receiptNumber"`
		Type             string      `json:"type"`
		TaxBaseBasic     float64     `json:"taxBaseBasic"`
		TaxBaseReduced   float64     `json:"taxBaseReduced"`
		TotalPrice       float64     `json:"totalPrice"`
		FreeTaxAmount    float64     `json:"freeTaxAmount"`
		VatAmountBasic   float64     `json:"vatAmountBasic"`
		VatAmountReduced float64     `json:"vatAmountReduced"`
		VatRateBasic     float64     `json:"vatRateBasic"`
		VatRateReduced   float64     `json:"vatRateReduced"`
		Items            []struct {
			Name     string  `json:"name"`
			ItemType string  `json:"itemType"`
			Quantity float64 `json:"quantity"`
			VatRate  float64 `json:"vatRate"`
			Price    float64 `json:"price"`
		} `json:"items"`
		Organization struct {
			BuildingNumber             interface{} `json:"buildingNumber"`
			Country                    string      `json:"country"`
			Dic                        string      `json:"dic"`
			IcDph                      string      `json:"icDph"`
			Ico                        string      `json:"ico"`
			Municipality               string      `json:"municipality"`
			Name                       string      `json:"name"`
			PostalCode                 string      `json:"postalCode"`
			PropertyRegistrationNumber string      `json:"propertyRegistrationNumber"`
			StreetName                 string      `json:"streetName"`
			VatPayer                   bool        `json:"vatPayer"`
		} `json:"organization"`
		Unit struct {
			CashRegisterCode           string      `json:"cashRegisterCode"`
			BuildingNumber             interface{} `json:"buildingNumber"`
			Country                    string      `json:"country"`
			Municipality               string      `json:"municipality"`
			PostalCode                 string      `json:"postalCode"`
			PropertyRegistrationNumber string      `json:"propertyRegistrationNumber"`
			StreetName                 string      `json:"streetName"`
			Name                       interface{} `json:"name"`
			UnitType                   string      `json:"unitType"`
		} `json:"unit"`
		Exemption bool `json:"exemption"`
	} `json:"receipt"`
	SearchIdentification struct {
		CreateDate        int64  `json:"createDate"`
		Bucket            int    `json:"bucket"`
		InternalReceiptID string `json:"internalReceiptId"`
		SearchUUID        string `json:"searchUuid"`
	} `json:"searchIdentification"`
}

type FinancnaSpravaBasic struct {
	ReceiptID string `json:"receiptId"`
	Items     []struct {
		Name     string  `json:"name"`
		ItemType string  `json:"itemType"`
		Quantity float64 `json:"quantity"`
		VatRate  float64 `json:"vatRate"`
		Price    float64 `json:"price"`
	} `json:"items"`
}
