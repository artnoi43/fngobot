package common

import "github.com/artnoi43/fngobot/internal/enums"

// Quote struct for Bitkub
type ApiQuote struct {
	Src  enums.Src
	Last float64
	Bid  float64
	Ask  float64
}

func (q *ApiQuote) QuoteLast() (float64, error) {
	return q.Last, nil
}
func (q *ApiQuote) QuoteBid() (float64, error) {
	return q.Bid, nil
}
func (q *ApiQuote) QuoteAsk() (float64, error) {
	return q.Ask, nil
}
