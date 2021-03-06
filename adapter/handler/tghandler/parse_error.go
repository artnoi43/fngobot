package tghandler

import (
	"strings"

	"github.com/artnoi43/fngobot/adapter/parse"
)

// HandleParsingError handles errors from package parse
func (h *handler) HandleParsingError(e parse.ParseError) {
	h.reply(formString(e))
}

func formString(e parse.ParseError) string {
	signals := []string{
		"failed to parse command:",
		parse.ErrMsgs[e].Error(),
	}
	return strings.Join(signals, "\n")
}
