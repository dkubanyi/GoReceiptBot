package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"time"
)

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
	FilePath         string              `json:"filePath"`
}

func GetReceiptsForUser(userId uuid.UUID) ([]Receipt, error) {
	var receipts []Receipt

	rows, err := DB.Query(`SELECT * FROM receipts WHERE id IN (SELECT receipt_id::uuid as uuid FROM user_receipts WHERE user_id::text = $1)`, userId)
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
			&r.FilePath,
		)

		receipts = append(receipts, r)
	}

	return receipts, err
}

func GetReceiptByReceiptId(receiptId string) (Receipt, error) {
	var r Receipt

	row := DB.QueryRow(`SELECT * FROM receipts WHERE receipt_id = $1`, receiptId)
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
		&r.FilePath,
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

func CreateReceipt(r Receipt, filePath string) (Receipt, error) {
	var receiptId string
	sqlStatement := "INSERT INTO receipts (id, receipt_id, cash_register_code, ico, ic_dph, dic, type, invoice_number, receipt_number, total_price, tax_base_basic, tax_base_reduced, vat_amount_basic, vat_amount_reduced, vat_rate_basic, vat_rate_reduced, issue_date, create_date, organization, unit, items, file_path)" +
		" VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22) RETURNING receipt_id;"

	layout := "02.01.2006 15:04:05"
	issueDate, _ := time.Parse(layout, r.IssueDate)
	createDate, _ := time.Parse(layout, r.CreateDate)

	org, _ := json.Marshal(r.Organization)
	unit, _ := json.Marshal(r.Unit)
	items, _ := json.Marshal(r.Items)

	err := DB.QueryRow(
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
		filePath,
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

	_, err := DB.Query("INSERT INTO user_receipts (user_id, receipt_id) VALUES ($1, $2);", u.Id, r.Id)
	if err != nil {
		return errors.New("failed to execute query")
	}

	return nil
}

func DeleteReceiptsByUserId(u *User) error {
	if u.Id == uuid.Nil {
		return errors.New("could not delete receipts of user: nil user ID passed")
	}

	receipts, err := GetReceiptsForUser(u.Id)
	if err != nil {
		return errors.New("could not fetch existing receipts for user")
	}

	for _, r := range receipts {
		if r.FilePath != "" {
			if err := os.Remove(r.FilePath); err != nil {
				msg := "could not delete file from the filesystem"
				log.Printf("%s: %v", msg, err)
				return errors.New(msg)
			}
		}
	}

	if _, err := DB.Query("DELETE FROM receipts WHERE id IN (SELECT receipt_id::uuid as uuid FROM user_receipts WHERE user_id::text = $1)", u.Id); err != nil {
		msg := "failed to delete receipts"
		log.Printf("%s: %v", msg, err)
		return errors.New(msg)
	}

	if err != nil {
		return errors.New("could not delete receipts of user, query failed")
	}

	return nil
}
