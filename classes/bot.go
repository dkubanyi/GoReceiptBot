package classes

import "log"

type Bot struct {
	handlers *Handlers
	msg      chan responseMessage
	done     chan struct{}
}

func New(h *Handlers) *Bot {
	if h.Errored == nil {
		h.Errored = logErrorHandler
	}

	b := &Bot{
		handlers: h,
		msg:      make(chan responseMessage),
	}

	go b.processMessages()

	return b
}

type Handlers struct {
	Response ResponseHandler
	Errored  ErrorHandler
}

type ResponseHandler func(OutgoingMessage)

type ErrorHandler func(msg string, err error)

type responseMessage struct {
	target, message string
}

func logErrorHandler(msg string, err error) {
	log.Printf("%s: %s", msg, err.Error())
}

func (b *Bot) processMessages() {
	for {
		select {
		case msg := <-b.msg:
			b.sendResponse(msg)
		case <-b.done:
			return
		}
	}
}

func (b *Bot) sendResponse(resp responseMessage) {
	b.handlers.Response(OutgoingMessage{
		Target:  resp.target,
		Message: resp.message,
	})
}

type OutgoingMessage struct {
	Target  string
	Message string
}
