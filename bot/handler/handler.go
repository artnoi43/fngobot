package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/artnoi43/fngobot/parse"
	"github.com/google/uuid"
	tb "gopkg.in/tucnak/telebot.v2"
)

// Bot types - used in Handler.Handle()
const (
	QUOTEBOT = iota
	TRACKBOT
	ALERTBOT
)

type Handler struct {
	uuid string
	quit chan bool
	conf Config
	cmd  *parse.BotCommand
	bot  *tb.Bot
	msg  *tb.Message
}

func NewHandler(b *tb.Bot, m *tb.Message, conf Config, cmd *parse.BotCommand) *Handler {
	uuid := strings.Split(uuid.NewString(), "-")[0]
	quit := make(chan bool, 1)
	log.Printf("[%s]: %s (from %d)\n", uuid, m.Text, m.Sender.ID)
	return &Handler{
		uuid: uuid,
		quit: quit,
		conf: conf,
		cmd:  cmd,
		bot:  b,
		msg:  m,
	}
}

func (h *Handler) Handle(t int) {
	switch t {
	case QUOTEBOT:
		h.SendQuote(h.cmd.Quote.Securities)
	case TRACKBOT:
		h.Track(h.cmd.Track.Securities, h.cmd.Track.TrackTimes, h.conf)
	case ALERTBOT:
		h.PriceAlert(h.cmd.Alert, h.conf)
	}
}

func (h *Handler) send(s string) {
	h.bot.Send(h.msg.Sender, s)
}

func (h *Handler) notifyStop() {
	log.Printf("[%s]: Received stop signal", h.uuid)
	h.send(fmt.Sprintf("Stopping %s", h.uuid))
}
