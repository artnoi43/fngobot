package clihandler

import (
	"log"
	"time"

	_handler "github.com/artnoi43/fngobot/adapter/handler"
	"github.com/artnoi43/fngobot/adapter/handler/utils"
	"github.com/artnoi43/fngobot/adapter/parse"
	"github.com/artnoi43/fngobot/internal/enums"
)

type handler struct {
	*_handler.BaseHandler
	conf *Config       `json:"-" yaml:"-"`
	done chan struct{} `json:"-" yaml:"-"`
}

func New(
	cmd *parse.BotCommand,
	conf *Config,
	done chan struct{},
) _handler.Handler {
	return &handler{
		BaseHandler: &_handler.BaseHandler{
			Start: time.Now(),
			Uuid:  utils.NewUUID(true),
			Cmd:   cmd,
			Quit:  utils.NewQuit(),
		},
		conf: conf,
		done: done,
	}
}

func (h *handler) UUID() string              { return h.Uuid }
func (h *handler) QuitChan() chan struct{}   { return h.Quit }
func (h *handler) GetCmd() *parse.BotCommand { return h.Cmd }
func (h *handler) Done()                     { h.IsDone = true }
func (h *handler) IsRunning() bool           { return !h.IsDone }

func (h *handler) Handle(t enums.BotType) {
	switch t {
	case enums.QuoteBot:
		h.Quote(
			h.GetCmd().Quote.Securities,
		)
	case enums.TrackBot:
		h.Track(
			h.GetCmd().Track.Securities,
			h.GetCmd().Track.TrackTimes,
		)
	case enums.AlertBot:
		h.PriceAlert(
			h.GetCmd().Alert,
		)
	}
}

func (h *handler) HandleParsingError(e parse.ParseError) {
	log.Printf(
		"[error] %s\n",
		parse.ErrMsgs[e],
	)
}

func (h *handler) notifyStop() {
	log.Printf("Stopping %s\n", h.UUID())
}

func (h *handler) StartTime() time.Time {
	return h.Start
}
