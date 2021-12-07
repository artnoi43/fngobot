package parse

import (
	"strconv"
	"strings"

	"github.com/artnoi43/fngobot/bot"
	"github.com/artnoi43/fngobot/enums"
	"github.com/artnoi43/fngobot/help"
)

const (
	ErrParseInt = iota + 1
	ErrParseFloat
	ErrInvalidSign
	ErrInvalidBidAskSwitch
	ErrInvalidQuoteTypeBid
	ErrInvalidQuoteTypeAsk
	ErrInvalidQuoteTypeLast
)

const (
	HelpCmd = iota
	QuoteCmd
	TrackCmd
	AlertCmd
)

// UserCommand is essentially a chat message.
// UserCommand.Chat is an int enum
type UserCommand struct {
	Command int
	Chat    string
}

type helpCommand struct {
	HelpMessage string
}

type quoteCommand struct {
	Securities []bot.Security
}

type trackCommand struct {
	quoteCommand
	TrackTimes int
}

// BotCommand is derived from UserCommand by parsing with Parse()
// Alerting does not need its own command struct,
// as the bot.Alert struct already has all the info needed.
type BotCommand struct {
	Help  helpCommand
	Quote quoteCommand
	Track trackCommand
	Alert bot.Alert
}

func getSrc(sw string) (idx int, src int) {
	switch sw {
	case "CRYPTO":
		idx = 2
		src = enums.YahooCrypto
	case "SATANG":
		idx = 2
		src = enums.Satang
	case "BITKUB":
		idx = 2
		src = enums.Bitkub
	default:
		idx = 1
		src = enums.Yahoo
	}
	return idx, src
}

func (cmd *quoteCommand) appendSecurities(ticks []string, src int) {
	for _, tick := range ticks {
		var s bot.Security
		s.Tick = strings.ToUpper(tick)
		s.Src = src
		cmd.Securities = append(cmd.Securities, s)
	}
}

// Parse parses UserCommand to BotCommand
func (c UserCommand) Parse() (cmd BotCommand, parseError int) {
	chat := strings.Split(c.Chat, " ")
	lenChat := len(chat)
	sw := strings.ToUpper(chat[1])
	idx, src := getSrc(sw)

	switch c.Command {
	case HelpCmd:
		cmd.Help.HelpMessage = help.GetHelp(c.Chat)
	case QuoteCmd:
		cmd.Quote.appendSecurities(chat[idx:], src)
	case TrackCmd:
		cmd.Track.appendSecurities(chat[idx:lenChat-1], src)
		r, err := strconv.Atoi(chat[lenChat-1])
		if err != nil {
			return cmd, ErrParseInt
		}
		cmd.Track.TrackTimes = r
	case AlertCmd:
		cmd.Alert.Security.Tick = strings.ToUpper(chat[idx])
		cmd.Alert.Security.Src = src
		targ, err := strconv.ParseFloat(chat[lenChat-1], 64)
		if err != nil {
			return cmd, ErrParseFloat
		}
		cmd.Alert.Target = targ

		sign := chat[lenChat-2]
		if sign != "<" && sign != ">" {
			return cmd, ErrInvalidSign
		} else if sign == "<" {
			cmd.Alert.Condition = enums.Lt
		} else if sign == ">" {
			cmd.Alert.Condition = enums.Gt
		}

		// Determine if bid/ask or last alert
		// bid/ask alert will have length = default length + 1
		defLen := make(map[int]int)
		defLen[enums.YahooCrypto] = 5 /* /alert crypto btc-usd > 1 */
		defLen[enums.Satang] = 5      /* /alert satang btc bid > 1 */
		defLen[enums.Bitkub] = 5      /* /alert bitkub btc > 1 */
		defLen[enums.Yahoo] = 4       /* /alert bbl.bk > 120 */
		// If unsupported alert/ use ones that are supported
		switch lenChat {
		// Last price alerts
		case defLen[cmd.Alert.Src]:
			// Satang does not support last price
			if cmd.Alert.Src == enums.Satang {
				return cmd, ErrInvalidQuoteTypeLast
			}
			cmd.Alert.QuoteType = enums.Last
		// Bid/ask alerts
		default:
			ba := strings.ToUpper(chat[idx+1])
			switch cmd.Alert.Src {
			// Yahoo Crypto does not support bid/ask price
			case enums.YahooCrypto:
				if ba == "BID" {
					return cmd, ErrInvalidQuoteTypeBid
				} else if ba == "ASK" {
					return cmd, ErrInvalidQuoteTypeAsk
				}
			default:
				if ba == "BID" {
					cmd.Alert.QuoteType = enums.Bid
				} else if ba == "ASK" {
					cmd.Alert.QuoteType = enums.Ask
				} else {
					return cmd, ErrInvalidBidAskSwitch
				}
			}
		}
	}
	return cmd, parseError
}
