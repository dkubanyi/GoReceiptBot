package handlers

import (
	"GoBudgetBot/constants"
	"GoBudgetBot/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	context models.BotContext
}

func (h *imageHandler) IsResponsible() bool {
	return len(h.context.Message.Photo) != 0
}

func (h *imageHandler) Process() error {
	// The last imgFile in h.imgFile slice has the best quality
	img := h.context.Message.Photo
	fileId := img[len(img)-1].FileID

	url := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", os.Getenv(constants.TelegramToken), fileId)

	resp, err := http.Get(url)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to fetch photo from url %s. Reason: %v", url, err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to read response body from Financna sprava. Reason: %v", err))
	}

	tgResponse := new(models.TelegramResponse)

	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&tgResponse); err != nil {
		log.Print(err)
		return errors.New("failed to decode QR code")
	}

	photoUrl := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", os.Getenv(constants.TelegramToken), tgResponse.Result.FilePath)

	fileName := strings.SplitAfter(photoUrl, "/")
	p := fmt.Sprintf(
		"%s/%s/%s/%s",
		"data",
		h.context.User.UserId,
		h.context.User.ChatId,
		fileName[len(fileName)-1],
	)

	imgFile, _ := createFileWithSubdirectories(p)
	defer imgFile.Close()

	r, _ := http.Get(photoUrl)
	defer resp.Body.Close()

	io.Copy(imgFile, r.Body)

	n, err := filepath.Abs(imgFile.Name())

	file, err := recognizeFile(n)

	if err != nil {
		return err
	}

	existingReceipt, err := models.GetReceiptByReceiptId(file.Receipt.ReceiptId)
	if existingReceipt.Id != uuid.Nil {
		return errors.New("this receipt already exists in the database")
	}

	filePath, err := filepath.Abs(imgFile.Name())
	if err != nil {
		return errors.New("failed to get absolute path from saved file")
	}

	// TODO transaction
	receipt, err := models.CreateReceipt(file.Receipt, filePath)
	if err != nil {
		log.Printf("could not create receipt: %v", err)
		return errors.New("failed to save receipt, please try again later")
	}

	if err := models.CreateUserReceiptMapping(h.context.User, &receipt); err != nil {
		log.Printf("could not create user-receipt mapping: %v", err)
		return errors.New("failed to save receipt mapping to user")
	}

	response = "Your receipt contains the following items:\n"
	for _, item := range file.Receipt.Items {
		response += fmt.Sprintf("<b>Item</b>: %s\n<b>Item type</b>: %s\n<b>Quantity</b>: %d pcs\n<b>VAT</b>: %d\n<b>Price</b>: %f\n\n", item.Name, item.ItemType, int64(item.Quantity), int64(item.VatRate), item.Price)
	}
	response += "\n That's all 😊"

	return nil
}

func (h *imageHandler) GetResponseMessage() string {
	if len(response) == 0 {
		return fmt.Sprintf("No QR code detected on the image. Please make sure it is well visible on the uploaded photo, and try again.")
	} else {
		return response
	}
}

func recognizeFile(path string) (*models.FinancnaSpravaResponse, error) {
	// open and decode image file
	file, _ := os.Open(path)
	img, _, err := image.Decode(file)

	if err != nil {
		log.Printf("Could not decode image: %v", err)
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

func verifyReceipt(receiptCode string) (*models.FinancnaSpravaResponse, error) {
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

	var response models.FinancnaSpravaResponse
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
