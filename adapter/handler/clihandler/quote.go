package clihandler

import (
	"fmt"
	"log"
	"sync"

	"github.com/artnoi43/fngobot/adapter/handler/utils"
	"github.com/artnoi43/fngobot/entity"
	"github.com/artnoi43/fngobot/internal/enums"
	"github.com/artnoi43/fngobot/usecase"
)

func (h *handler) Quote(securities []usecase.Security) {
	var wg sync.WaitGroup
	for _, security := range securities {
		wg.Add(1)
		go func(s usecase.Security) {
			defer wg.Done()
			q, err := s.Quote()
			if err != nil {
				log.Printf(
					"Failed to fetch %s quote from %s: %s\n",
					s.Tick,
					s.Src.String(),
					err.Error(),
				)
				return
			}
			printQuote(s.Tick, s.Src, q)
		}(security)
	}
	wg.Wait()
}

func printQuote(t string, s enums.Src, q entity.Quoter) {
	bid, err := q.QuoteBid()
	if err != nil {
		bid = -1
	}
	ask, err := q.QuoteAsk()
	if err != nil {
		ask = -1
	}
	last, err := q.QuoteLast()
	if err != nil {
		last = -1
	}
	utils.Printer.Printf(
		"Ticker: %s [%s]\nBid: %f\nAsk: %f\nLast: %f\n",
		t, s, bid, ask, last,
	)
	fmt.Println(enums.Bar)
}
