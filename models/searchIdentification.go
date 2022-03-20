package models

type SearchIdentification struct {
	CreateDate        int64  `json:"createDate"`
	Bucket            int    `json:"bucket"`
	InternalReceiptID string `json:"internalReceiptId"`
	SearchUUID        string `json:"searchUuid"`
}
