package models

// FinancnaSpravaResponse
// Following structs represent the response from Financna sprava API. Official document is available here:
// https://www.financnasprava.sk/_img/pfsedit/Dokumenty_PFS/Podnikatelia/eKasa/2019/2019.05.27_eKasa_rozhranie.pdf
type FinancnaSpravaResponse struct {
	ReturnValue          int                  `json:"returnValue"`
	Receipt              Receipt              `json:"receipt"`
	SearchIdentification SearchIdentification `json:"searchIdentification"`
}
