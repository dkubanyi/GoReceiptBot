package handlers

import (
	"GoBudgetBot/internal/domain/context"
	r "GoBudgetBot/internal/domain/receipt"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"strings"
)

/**
* This handler is responsible for processing of QR codes, passed as a string
 */
type qrHandler struct {
	context context.BotContext
}

func (h *qrHandler) IsResponsible() bool {
	return strings.HasPrefix(h.context.Message.Text, CommandQr)
}

func (h *qrHandler) Process() error {
	qr := strings.SplitAfter(h.context.Message.Text, " ")
	if len(qr) != 2 {
		return errors.New("wrong command. The /qr command needs to be followed by a space, and then the code. Example usage: /qr abcdef")
	}

	receipt, err := verifyReceipt(qr[1])
	if err != nil {
		return err
	}

	existingReceipt, err := r.GetReceiptByReceiptId(receipt.Receipt.ReceiptId)
	if existingReceipt.Id != uuid.Nil {
		return errors.New("this receipt already exists in the database")
	}

	// TODO transaction
	r, err := r.CreateReceipt(receipt.Receipt, "")
	if err != nil {
		log.Printf("could not create receipt: %v", err)
		return errors.New("failed to save receipt, please try again later")
	}

	if err := r.CreateUserReceiptMapping(h.context.User); err != nil {
		log.Printf("could not create user-receipt mapping: %v", err)
		return errors.New("failed to save receipt mapping to user")
	}

	response = "Your receipt contains the following items:\n"
	for _, item := range receipt.Receipt.Items {
		response += fmt.Sprintf("<b>Item</b>: %s\n<b>Item type</b>: %s\n<b>Quantity</b>: %d pcs\n<b>VAT</b>: %d\n<b>Price</b>: %f\n\n", item.Name, item.ItemType, int64(item.Quantity), int64(item.VatRate), item.Price)
	}
	response += "\n That's all 😊"

	return nil
}

func (h *qrHandler) GetResponseMessage() string {
	return response
}
