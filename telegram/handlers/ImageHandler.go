package handlers

import (
	"GoBudgetBot/constants"
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

/**
* This handler is responsible for processing updates containing images
 */
type imageHandler struct {
	text     string
	image    []tgbotapi.PhotoSize
	photoUrl string
}

func (h *imageHandler) IsResponsible() bool {
	return len(h.image) != 0
}

func (h *imageHandler) Process() {
	fileId := h.image[0].FileID
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

	photoUrl := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", os.Getenv(constants.TelegramToken), tgResponse.Result.FilePath)
	h.photoUrl = photoUrl
}

func (h *imageHandler) GetResponseMessage() string {
	// TODO process the photo, e.g. if it is a QR code of a payment, parse it and save in DB
	return fmt.Sprintf("Your photo is here: %s", h.photoUrl)
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
