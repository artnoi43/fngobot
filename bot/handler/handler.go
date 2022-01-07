package handler

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/artnoi43/fngobot/bot"
	"github.com/artnoi43/fngobot/parse"
	"github.com/google/uuid"
	tb "gopkg.in/tucnak/telebot.v3"
)

// Bot types - used in handler.Handle()
const (
	QUOTEBOT = iota
	TRACKBOT
	ALERTBOT
	HANDLERS
)

type Handler interface {
	// These exported methods are called from other packages
	UUID() string
	QuitChan() chan bool
	Done()
	GetCmd() *parse.BotCommand
	Handle(int)
	HandleParsingError(parse.ParseError)
	SendQuote([]bot.Security)
	Track([]bot.Security, int, Config)
	PriceAlert(bot.Alert, Config)
	// These unexported methods are called from this package
	send(string)
	yaml() string
	isRunning() bool
}

type Handlers []Handler

var BotHandlers Handlers

type handler struct {
	Uuid   string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	Cmd    *parse.BotCommand `json:"command,omitempty" yaml:"command,omitempty"`
	Start  time.Time         `json:"start,omitempty" yaml:"start,omitempty"`
	Quit   chan bool         `json:"-" yaml:"-"`
	IsDone bool              `json:"-" yaml:"-"`
	Conf   Config            `json:"-" yaml:"-"`
	Bot    *tb.Bot           `json:"-" yaml:"-"`
	Msg    *tb.Message       `json:"-" yaml:"-"`
}

func (h *handler) UUID() string {
	return h.Uuid
}
func (h *handler) QuitChan() chan bool {
	return h.Quit
}
func (h *handler) Done() {
	h.IsDone = true
}
func (h *handler) isRunning() bool {
	return !h.IsDone
}
func (h *handler) GetCmd() *parse.BotCommand {
	return h.Cmd
}

// Handle calls different methods on h based on its function parameter
func (h *handler) Handle(t int) {
	switch t {
	case QUOTEBOT:
		h.SendQuote(h.Cmd.Quote.Securities)
	case TRACKBOT:
		h.Track(h.Cmd.Track.Securities, h.Cmd.Track.TrackTimes, h.Conf)
	case ALERTBOT:
		h.PriceAlert(h.Cmd.Alert, h.Conf)
	case HANDLERS:
		h.SendHandlers()
	}
}

func (h *handler) send(s string) {
	h.Bot.Send(h.Msg.Sender, s)
}

func (h *handler) notifyStop() {
	log.Printf("[%s]: Received stop signal", h.Uuid)
	h.send(fmt.Sprintf("Stopping %s", h.Uuid))
}

// Stop stops a handler with matching UUID
func (h *Handlers) Stop(uuid string) (i int, ok bool) {
	for idx, handler := range *h {
		switch uuid {
		case handler.UUID():
			log.Printf("[%s]: Sending quit signal\n", handler.UUID())
			quit := handler.QuitChan()
			quit <- true
			log.Printf("[%s]: Sent quit signal\n", handler.UUID())
			i = idx
			ok = true
		}
	}
	return i, ok
}

// NewHandler returns a new handler and appends it to BotHandlers
func NewHandler(b *tb.Bot, m *tb.Message, conf Config, cmd *parse.BotCommand) Handler {
	uuid := strings.Split(uuid.NewString(), "-")[0]
	quit := make(chan bool, 1)
	log.Printf("[%s]: %s (from %d)\n", uuid, m.Text, m.Sender.ID)
	h := &handler{
		Uuid:   uuid,
		Start:  time.Now(),
		Cmd:    cmd,
		Quit:   quit,
		IsDone: false,
		Conf:   conf,
		Bot:    b,
		Msg:    m,
	}
	BotHandlers = append(BotHandlers, h)
	return h
}
