package tghandler

import (
	"errors"
	"sync"

	"github.com/artnoi43/fngobot/bot"
	"github.com/artnoi43/fngobot/bot/utils"
	"github.com/artnoi43/fngobot/fetch"
)

// SendQuote sends quote(s) to users via chats.
// It is reused by tracking and alerting handlers.
func (h *handler) Quote(securities []bot.Security) {
	quotes := make(chan fetch.Quoter, len(securities))
	var wg sync.WaitGroup
	for _, security := range securities {
		wg.Add(1)
		// This Goroutines get quotes
		go func() {
			defer wg.Done()
			q, err := security.Quote()
			if err != nil {
				var errMsg string
				if errors.Is(err, fetch.ErrNotFound) {
					errMsg = "Ticker not found"
				} else {
					errMsg = "Error getting quote"
				}
				h.send(utils.Printer.Sprintf(
					"[%s]\n%s: %s from %s",
					h.UUID(),
					errMsg,
					security.Tick,
					security.GetSrcStr(),
				))
				return
			}
			quotes <- q
		}()

		wg.Add(1)
		// This Goroutine sends quotes to users
		go func(s bot.Security) {
			defer wg.Done()
			for q := range quotes {
				last, _ := q.Last()
				bid, _ := q.Bid()
				ask, _ := q.Ask()
				msg := utils.Printer.Sprintf(
					"[%s]\nQuote from %s\n%s\nBid: %f\nAsk: %f\nLast: %f\n",
					h.UUID(),
					s.GetSrcStr(),
					s.Tick,
					bid,
					ask,
					last,
				)
				h.send(msg)
			}
		}(security)
	}
	wg.Wait()
}
