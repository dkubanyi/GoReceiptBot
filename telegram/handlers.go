package telegram

import "errors"

type ResponseHandler interface {
	isResponsible() bool
	process()
	getResponseMessage() string
}

func InitHandler(msg string) (ResponseHandler, error) {
	handlers := []ResponseHandler{
		startHandler{msgReceived: msg},
		testHandler{msgReceived: msg},
	}

	for _, h := range handlers {
		if h.isResponsible() {
			return h, nil
		}
	}

	return nil, errors.New("handler not implemented")
}

type startHandler struct {
	msgReceived string
}

func (h startHandler) isResponsible() bool {
	return h.msgReceived == "/start"
}

func (h startHandler) process() {
	// TODO save into DB, etc
}

func (h startHandler) getResponseMessage() string {
	return defaultMessage
}

type testHandler struct {
	msgReceived string
}

func (h testHandler) isResponsible() bool {
	return h.msgReceived == "/end"
}

func (h testHandler) process() {
	// TODO save into DB, etc
}

func (h testHandler) getResponseMessage() string {
	return "bye bye!"
}
