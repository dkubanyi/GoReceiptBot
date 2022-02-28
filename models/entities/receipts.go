package entities

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

// FinancnaSpravaResponse
// Following structs represent the response from Financna sprava. Official document is available here:
// https://www.financnasprava.sk/_img/pfsedit/Dokumenty_PFS/Podnikatelia/eKasa/2019/2019.05.27_eKasa_rozhranie.pdf
type FinancnaSpravaResponse struct {
	ReturnValue          int                  `json:"returnValue"`
	Receipt              Receipt              `json:"receipt"`
	SearchIdentification SearchIdentification `json:"searchIdentification"`
}

type SearchIdentification struct {
	CreateDate        int64  `json:"createDate"`
	Bucket            int    `json:"bucket"`
	InternalReceiptID string `json:"internalReceiptId"`
	SearchUUID        string `json:"searchUuid"`
}

type Receipt struct {
	Id               uuid.UUID           `json:"id"`
	ReceiptId        string              `json:"receiptId"`
	CashRegisterCode string              `json:"cashRegisterCode"`
	Ico              string              `json:"ico"`
	IcDph            string              `json:"icDph"`
	Dic              string              `json:"dic"`
	Type             string              `json:"type"`
	InvoiceNumber    string              `json:"invoiceNumber"`
	ReceiptNumber    int64               `json:"receiptNumber"`
	TotalPrice       float64             `json:"totalPrice"`
	TaxBaseBasic     float64             `json:"taxBaseBasic"`
	TaxBaseReduced   float64             `json:"taxBaseReduced"`
	VatAmountBasic   float64             `json:"vatAmountBasic"`
	VatAmountReduced float64             `json:"vatAmountReduced"`
	VatRateBasic     float64             `json:"vatRateBasic"`
	VatRateReduced   float64             `json:"vatRateReduced"`
	Exemption        bool                `json:"exemption"`
	IssueDate        string              `json:"issueDate"`
	CreateDate       string              `json:"createDate"`
	Organization     ReceiptOrganization `json:"organization"`
	Unit             ReceiptUnit         `json:"unit"`
	Items            []ReceiptItem       `json:"items"`
}

type ReceiptItem struct {
	Name     string  `json:"name"`
	ItemType string  `json:"itemType"`
	Quantity float64 `json:"quantity"`
	VatRate  float64 `json:"vatRate"`
	Price    float64 `json:"price"`
}

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

func GetReceiptByReceiptId(receiptId string) (Receipt, error) {
	db := CreateConnection()
	defer db.Close()
	var r Receipt

	row := db.QueryRow(`SELECT * FROM receipts WHERE receipt_id = $1`, receiptId)

	err := row.Scan(
		&r.Id,
		&r.ReceiptId,
		&r.CashRegisterCode,
		&r.Ico,
		&r.IcDph,
		&r.Dic,
		&r.Type,
		&r.InvoiceNumber,
		&r.ReceiptNumber,
		&r.TotalPrice,
		&r.TaxBaseBasic,
		&r.TaxBaseReduced,
		&r.VatAmountBasic,
		&r.VatAmountReduced,
		&r.VatRateBasic,
		&r.VatRateReduced,
		&r.Exemption,
		&r.IssueDate,
		&r.CreateDate,
		&r.Unit,
		"{}", //&r.Organization, //&r.Organization, // convert from json
		"{}", //&r.Items,        //&r.Items,        // convert from json
	)

	switch err {
	case sql.ErrNoRows:
		return r, errors.New("no rows were returned")
	case nil:
		return r, nil
	default:
		return r, errors.New(fmt.Sprintf("unable to scan the row. %v", err))
	}
}

func CreateReceipt(r Receipt) (Receipt, error) {
	db := CreateConnection()
	defer db.Close()

	var receiptId string
	sqlStatement := "INSERT INTO receipts (id, receipt_id, cash_register_code, ico, icdph, dic, type, invoice_number, receipt_number, total_price, tax_base_basic, tax_base_reduced, vat_amount_basic, vat_amount_reduced, vat_rate_basic, vat_rate_reduced, exemption, issue_date, create_date, organization, unit, items)" +
		" VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22) RETURNING receipt_id;"

	layout := "02.01.2006 15:04:05"
	issueDate, _ := time.Parse(layout, r.IssueDate)
	createDate, _ := time.Parse(layout, r.CreateDate)

	org, _ := json.Marshal(r.Organization)
	unit, _ := json.Marshal(r.Unit)
	items, _ := json.Marshal(r.Items)

	err := db.QueryRow(
		sqlStatement,
		uuid.New(),
		r.ReceiptId,
		r.CashRegisterCode,
		r.Ico,
		r.IcDph,
		r.Dic,
		r.Type,
		r.InvoiceNumber,
		r.ReceiptNumber,
		r.TotalPrice,
		r.TaxBaseBasic,
		r.TaxBaseReduced,
		r.VatAmountBasic,
		r.VatAmountReduced,
		r.VatRateBasic,
		r.VatRateReduced,
		r.Exemption,
		issueDate,
		createDate,
		org,
		unit,
		items,
	).Scan(&receiptId)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	log.Printf("Inserted new receipt with receipt_id %v", receiptId)
	return GetReceiptByReceiptId(receiptId)
}
