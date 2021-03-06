package fetch

import (
	"log"

	"github.com/artnoi43/fngobot/adapter/fetch/binance"
	"github.com/artnoi43/fngobot/adapter/fetch/bitkub"
	"github.com/artnoi43/fngobot/adapter/fetch/coinbase"
	yahoocrypto "github.com/artnoi43/fngobot/adapter/fetch/crypto"
	"github.com/artnoi43/fngobot/adapter/fetch/satang"
	"github.com/artnoi43/fngobot/adapter/fetch/yahoo"
	"github.com/artnoi43/fngobot/internal/enums"
	"github.com/artnoi43/fngobot/usecase"
)

type newFunc func() interface{}

var (
	newFetcherMap = map[enums.Src]newFunc{
		enums.Yahoo:       yahoo.New,
		enums.YahooCrypto: yahoocrypto.New,
		enums.Satang:      satang.New,
		enums.Bitkub:      bitkub.New,
		enums.Binance:     binance.New,
		enums.Coinbase:    coinbase.New,
	}
)

func New(s enums.Src) usecase.Fetcher {
	f := newFetcherMap[s]
	fetcher, ok := f().(usecase.Fetcher)
	if !ok {
		log.Fatalf("fetcher %s not implementing usecase.Fetcher\n", s)
	}
	if fetcher == nil {
		log.Fatalf("nil fetcher for %s\n", s)
	}
	return fetcher
}
