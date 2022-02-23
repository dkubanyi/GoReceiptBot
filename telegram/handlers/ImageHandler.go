package handlers

import (
	"GoBudgetBot/constants"
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	tgResponse := new(TelegramResponse)
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&tgResponse)
	if err != nil {
		log.Fatal(err)
		return
	}

	h.photoUrl = fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", os.Getenv(constants.TelegramToken), tgResponse.Result.FilePath)

	// TODO save in structure: {dataFolder}/{userId}/{chatId}
	fileName := strings.SplitAfter(h.photoUrl, "/")
	img, _ := os.Create(fileName[len(fileName)-1])
	defer img.Close()

	r, _ := http.Get(h.photoUrl)
	defer resp.Body.Close()

	io.Copy(img, r.Body)

	n, err := filepath.Abs(img.Name())
	h.parsedQrString = recognizeFile(n)
}

func (h *imageHandler) GetResponseMessage() string {
	var msg string

	if len(h.parsedQrString) == 0 {
		msg = fmt.Sprintf("No QR code detected on the image. Please make sure it is well visible on the uploaded photo, and try again.")
	} else {
		msg = fmt.Sprintf("Parsed QR text: %s. URL of the photo on Telegram's servers: %s", h.parsedQrString, h.photoUrl)
	}

	// TODO process the photo, e.g. if it is a QR code of a payment, parse it and save in DB
	return msg
}

func recognizeFile(path string) string {

	// open and decode image file
	file, _ := os.Open(path)
	img, _, _ := image.Decode(file)

	if img == nil {
		log.Println("Could not decode image")
		return ""
	}

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)

	if err != nil {
		log.Println(err)
		return ""
	}

	return result.GetText()
}

type TelegramResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		FileID       string `json:"file_id"`
		FileUniqueID string `json:"file_unique_id"`
		FileSize     int    `json:"file_size"`
		FilePath     string `json:"file_path"`
	} `json:"result"`
}
