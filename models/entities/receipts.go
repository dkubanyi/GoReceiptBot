package entities

import (
	"database/sql"
	"database/sql/driver"
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

/**
Type:
- Pokladničný doklad - doklad o prijatí tržby, doklad o vrátení platby za vrátený tovar, doklad o
vrátení platby za tovar alebo poskytnutú službu pri ich reklamácii, doklad o vrátení platby za
neposkytnutú službu alebo doklad za vrátené zálohované obaly vyhotovený ORP
- Úhrada faktúry - doklad o prijatí tržby v prípade, ak ide o úhradu faktúry alebo jej časti
- Neplatný doklad - doklad vyhotovený v rámci skúšobnej prevádzky ORP alebo pri zaškolení
fyzickej osoby, ktorá bude evidovať tržby v ORP
- Vklad – doklad evidovaný pri vklade hotovosti do pokladnice
- Výber – doklad evidovaný pri výbere hotovosti z pokladnice
*/
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
	IssueDate        string              `json:"issueDate"`
	CreateDate       string              `json:"createDate"`
	Organization     ReceiptOrganization `json:"organization"`
	Unit             ReceiptUnit         `json:"unit"`
	Items            ReceiptItems        `json:"items"`
}

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

func GetReceiptsForUser(userId uuid.UUID) ([]Receipt, error) {
	db := CreateConnection()
	defer db.Close()
	var receipts []Receipt

	rows, err := db.Query(`SELECT * FROM receipts_sk WHERE id IN (SELECT receipt_id_sk::uuid as uuid FROM user_receipts WHERE user_id::text = $1)`, userId)
	defer rows.Close()

	if err != nil {
		log.Fatalf("Unable to execute query. %v", err)
	}

	for rows.Next() {
		var r Receipt
		err = rows.Scan(
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
			&r.IssueDate,
			&r.CreateDate,
			&r.Organization,
			&r.Unit,
			&r.Items,
		)

		receipts = append(receipts, r)
	}

	return receipts, err
}

func GetReceiptByReceiptId(receiptId string) (Receipt, error) {
	db := CreateConnection()
	defer db.Close()
	var r Receipt

	row := db.QueryRow(`SELECT * FROM receipts_sk WHERE receipt_id = $1`, receiptId)
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
		&r.IssueDate,
		&r.CreateDate,
		&r.Organization,
		&r.Unit,
		&r.Items,
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
	sqlStatement := "INSERT INTO receipts_sk (id, receipt_id, cash_register_code, ico, ic_dph, dic, type, invoice_number, receipt_number, total_price, tax_base_basic, tax_base_reduced, vat_amount_basic, vat_amount_reduced, vat_rate_basic, vat_rate_reduced, issue_date, create_date, organization, unit, items)" +
		" VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21) RETURNING receipt_id;"

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

func CreateUserReceiptMapping(u *User, r *Receipt) error {
	if u.Id == uuid.Nil || r.Id == uuid.Nil {
		return errors.New("could not create mapping for user and receipt")
	}

	db := CreateConnection()
	defer db.Close()

	_, err := db.Query("INSERT INTO user_receipts (id, user_id, receipt_id_sk) VALUES ($1, $2, $3);", uuid.New(), u.Id, r.Id)
	if err != nil {
		return errors.New("failed to execute query")
	}

	return nil
}

func DeleteReceiptsByUserId(u *User) error {
	if u.Id == uuid.Nil {
		return errors.New("could not delete receipts of user: nil user ID passed")
	}

	db := CreateConnection()
	defer db.Close()

	_, err := db.Query("DELETE FROM receipts_sk WHERE id IN (SELECT receipt_id_sk::uuid as uuid FROM user_receipts WHERE user_id::text = $1)", u.Id)

	// @TODO also delete saved photos from filesystem?

	if err != nil {
		return errors.New("could not delete receipts of user, query failed")
	}

	return nil
}
