package handlers

import (
	"GoBudgetBot/constants"
	"GoBudgetBot/models/entities"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

/**
* This handler is responsible for processing updates containing images
 */
type imageHandler struct {
	text           string
	image          []tgbotapi.PhotoSize
	photoUrl       string
	parsedQrString string
	user           entities.User
}

func (h *imageHandler) IsResponsible() bool {
	return len(h.image) != 0
}

func (h *imageHandler) Process() {
	// The last image in h.image slice has the best quality
	fileId := h.image[len(h.image)-1].FileID

	url := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", os.Getenv(constants.TelegramToken), fileId)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	tgResponse := new(entities.TelegramResponse)
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&tgResponse)
	if err != nil {
		log.Fatal(err)
		return
	}

	h.photoUrl = fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", os.Getenv(constants.TelegramToken), tgResponse.Result.FilePath)

	fileName := strings.SplitAfter(h.photoUrl, "/")
	p := fmt.Sprintf("%s/%s/%s/%s", "data", h.user.UserId, h.user.ChatId, fileName[len(fileName)-1])
	img, _ := createFileWithSubdirectories(p)
	defer img.Close()

	r, _ := http.Get(h.photoUrl)
	defer resp.Body.Close()

	io.Copy(img, r.Body)

	n, err := filepath.Abs(img.Name())

	file, err := recognizeFile(n)

	if err != nil {
		h.parsedQrString = err.Error()
		return
	}

	responseStr := "Your receipt contains the following items:\n"
	for _, item := range file.Receipt.Items {
		responseStr += fmt.Sprintf("<b>Item</b>: %s\n<b>Item type</b>: %s\n<b>Quantity</b>: %d pcs\n<b>VAT</b>: %d\n<b>Price</b>: %f\n\n", item.Name, item.ItemType, int64(item.Quantity), int64(item.VatRate), item.Price)
	}

	responseStr += "\n That's all 😊"

	// check if receipt_id already exists in db, if yes, return a message that it exists
	existingReceipt, err := entities.GetReceiptByReceiptId(file.Receipt.ReceiptId)
	if existingReceipt.Id != uuid.Nil {
		h.parsedQrString = "This receipt already exists in the database"
	} else {
		// if it doesn't yet exist in db, persist in db
		entities.CreateReceipt(file.Receipt)
		h.parsedQrString = responseStr
	}
}

func (h *imageHandler) GetResponseMessage() string {
	var msg string

	if len(h.parsedQrString) == 0 {
		msg = fmt.Sprintf("No QR code detected on the image. Please make sure it is well visible on the uploaded photo, and try again.")
	} else {
		msg = fmt.Sprintf("Result: %s. URL of the photo on Telegram's servers: %s", h.parsedQrString, h.photoUrl)
	}

	// TODO process the photo, e.g. if it is a QR code of a payment, parse it and save in DB
	return msg
}

func recognizeFile(path string) (*entities.FinancnaSpravaResponse, error) {
	// open and decode image file
	file, _ := os.Open(path)
	img, _, _ := image.Decode(file)

	if img == nil {
		log.Println("Could not decode image")
		return nil, errors.New("could not decode image")
	}

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)

	if err != nil {
		return nil, err
	}

	receipt, err := verifyReceipt(result.GetText())
	if err != nil {
		return nil, err
	}

	if receipt.Receipt.ReceiptId == "" {
		// not a valid receipt
		return nil, errors.New("receipt not recognized by Financna sprava")
	}

	return receipt, nil
}

func verifyReceipt(receiptCode string) (*entities.FinancnaSpravaResponse, error) {
	// request financna sprava
	finspravaUrl := "https://ekasa.financnasprava.sk/mdu/api/v1/opd/receipt/find"

	req, err := http.NewRequest("POST", finspravaUrl, bytes.NewBuffer([]byte(fmt.Sprintf(`{"receiptId": "%s"}`, receiptCode))))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Request to Financna sprava failed: %s", err))
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var response entities.FinancnaSpravaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New("cannot unmarshal JSON")
	}

	return &response, nil
}

func createFileWithSubdirectories(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}
