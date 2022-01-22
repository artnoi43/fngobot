package bot

import (
	"strings"

	"github.com/artnoi43/fngobot/enums"
)

// GetSrcStr returns source in string based on s.Src
func (s *Security) GetSrcStr() string {
	switch s.Src {
	case enums.Yahoo:
		return "Yahoo Finance"
	case enums.YahooCrypto:
		return "Yahoo Crypto"
	default:
		return strings.Title(string(s.Src))
	}
}

// GetCondStr returns condition in string based on a.Condition
func (a *Alert) GetCondStr() string {
	switch a.Condition {
	case enums.Lt:
		return "<="
	}
	return ">="
}

// GetQuoteTypeStr returns quote type in string based on a.QuoteType
func (a *Alert) GetQuoteTypeStr() string {
	switch a.QuoteType {
	case enums.Bid:
		return "BID"
	case enums.Ask:
		return "ASK"
	default:
		return "LAST"
	}
}
