package handlers

import (
	"GoBudgetBot/models"
	"GoBudgetBot/models/entities"
	"fmt"
)

/**
* This handler is responsible for responding to receipt queries
 */
type receiptHandler struct {
	context models.BotContext
}

func (h *receiptHandler) IsResponsible() bool {
	return h.context.Message.Text == CommandShowReceipts
}

func (h *receiptHandler) Process() error {
	return nil
}

func (h *receiptHandler) GetResponseMessage() string {
	receipts, err := entities.GetReceiptsForUser(h.context.User.Id)

	if err != nil {
		return fmt.Sprintf("An error occurred while fetching your data. Please try again later.\nDetails: %s", err.Error())
	}

	if len(receipts) == 0 {
		return "You have no saved receipts yet"
	}

	msg := "Your receipts:\n"
	for _, r := range receipts {
		msg += fmt.Sprintf("<b>Receipt ID:</b> %s\n<b>Issued</b> %s\n<b>Created</b> %s\n\n<b>Items:</b>\n", r.ReceiptId, r.IssueDate, r.CreateDate)

		for _, i := range r.Items {
			msg += fmt.Sprintf("<b>Name</b>: %s\n<b>Type</b>: %s\n<b>Qty</b>: %d pcs\n<b>VAT</b>: %d\n<b>Price</b>: %f€\n\n", i.Name, i.ItemType, int64(i.Quantity), int64(i.VatRate), i.Price)
		}
	}

	return msg
}
