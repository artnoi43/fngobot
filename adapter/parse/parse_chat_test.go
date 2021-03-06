package parse

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/artnoi43/fngobot/internal/enums"
	"github.com/artnoi43/fngobot/usecase"
)

func TestGetSrc(t *testing.T) {
	type getSrcOut struct {
		idx int
		src enums.Src
	}
	tests := []struct {
		In       string
		Expected getSrcOut
	}{

		{
			In: "YAHOO",
			Expected: getSrcOut{
				idx: 1,
				src: enums.Yahoo,
			},
		},
		{
			In: "CRYPTO",
			Expected: getSrcOut{
				idx: 2,
				src: enums.YahooCrypto,
			},
		},
		{
			In: "SATANG",
			Expected: getSrcOut{
				idx: 2,
				src: enums.Satang,
			},
		},
		{
			In: "BITKUB",
			Expected: getSrcOut{
				idx: 2,
				src: enums.Bitkub,
			},
		},
		{
			In: "BINANCE",
			Expected: getSrcOut{
				idx: 2,
				src: enums.Binance,
			},
		},
		{
			In: "COINBASE",
			Expected: getSrcOut{
				idx: 2,
				src: enums.Coinbase,
			},
		},
	}
	for _, test := range tests {
		src, idx := getSrc(test.In)
		if test.Expected.idx != idx {
			t.Errorf("invalid idx for %s\n", test.In)
		}
		if test.Expected.src != src {
			t.Errorf("invalid src for %s\n", test.In)
		}
	}
}

// Test parsing UserCommand into BotCommand
func TestParse(t *testing.T) {
	for _, test := range tests {
		in := test.In
		out, err := in.Parse()
		if err != 0 {
			t.Errorf("error parsing UserCommand: %+v\n", in)
		}

		report := func() {
			fmt.Printf("In: %s\nOut: %+v\n", in, out)
		}
		compareSecs := func() {
			for idx, actual := range out.Quote.Securities {
				expectedSecurity := test.Expected.Quote.Securities[idx]
				if actual.Src != expectedSecurity.Src {
					t.Errorf("[invalid security]: sources differ - actual: %s, expected: %s\n", actual.Src.String(), expectedSecurity.Src.String())
				}
				if actual.Tick != expectedSecurity.Tick {
					t.Errorf("[invalid security]: ticker symbols differ - actual: %s, expected: %s\n", actual.Tick, expectedSecurity.Tick)
				}
			}
		}
		switch test.In.TargetBot {
		case enums.QuoteBot:
			compareSecs()
		case enums.TrackBot:
			compareSecs()
			if out.Track.TrackTimes != test.Expected.Track.TrackTimes {
				t.Errorf("[/track] invalid track time\n")
				report()
			}
		case enums.AlertBot:
			if !reflect.DeepEqual(out.Alert, test.Expected.Alert) {
				compareSecs()
				if out.Alert.Src != test.Expected.Alert.Src {
					t.Errorf("[/alert] invalid alert source\n")
					report()
				}
				if out.Alert.Condition != test.Expected.Alert.Condition {
					t.Errorf("[/alert] invalid alert condition\n")
					report()
				}
				if out.Alert.QuoteType != test.Expected.Alert.QuoteType {
					t.Errorf("[/alert] invalid alert quote type\n")
					report()
				}
				if out.Alert.Target != test.Expected.Alert.Target {
					t.Errorf("invalid alert target\n")
					report()
				}
			}
		}
	}
}

type parseTest struct {
	In       UserCommand
	Expected BotCommand
}

var (
	tests = []parseTest{
		{
			In: UserCommand{
				Text:      "/quote gc=f",
				TargetBot: enums.QuoteBot,
			},
			Expected: BotCommand{
				Quote: quoteCommand{
					Securities: []usecase.Security{
						{
							Tick: "GC=F",
							Src:  enums.Yahoo,
						},
					},
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/quote satang btc",
				TargetBot: enums.QuoteBot,
			},
			Expected: BotCommand{
				Quote: quoteCommand{
					Securities: []usecase.Security{
						{
							Tick: "BTC",
							Src:  enums.Satang,
						},
					},
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/quote bitkub btc",
				TargetBot: enums.QuoteBot,
			},
			Expected: BotCommand{
				Quote: quoteCommand{
					Securities: []usecase.Security{
						{
							Tick: "BTC",
							Src:  enums.Bitkub,
						},
					},
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/quote binance btc",
				TargetBot: enums.QuoteBot,
			},
			Expected: BotCommand{
				Quote: quoteCommand{
					Securities: []usecase.Security{
						{
							Tick: "BTC",
							Src:  enums.Binance,
						},
					},
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/quote coinbase ada",
				TargetBot: enums.QuoteBot,
			},
			Expected: BotCommand{
				Quote: quoteCommand{
					Securities: []usecase.Security{
						{
							Tick: "ADA",
							Src:  enums.Coinbase,
						},
					},
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/track gc=f 2",
				TargetBot: enums.TrackBot,
			},
			Expected: BotCommand{
				Track: trackCommand{
					quoteCommand: quoteCommand{
						Securities: []usecase.Security{
							{
								Tick: "GC=F",
								Src:  enums.Yahoo,
							},
						},
					},
					TrackTimes: 2,
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/track satang btc 69",
				TargetBot: enums.TrackBot,
			},
			Expected: BotCommand{
				Track: trackCommand{
					quoteCommand: quoteCommand{
						Securities: []usecase.Security{
							{
								Tick: "BTC",
								Src:  enums.Satang,
							},
						},
					},
					TrackTimes: 69,
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/track bitkub btc ada 69",
				TargetBot: enums.TrackBot,
			},
			Expected: BotCommand{
				Track: trackCommand{
					quoteCommand: quoteCommand{
						Securities: []usecase.Security{
							{Tick: "BTC", Src: enums.Bitkub},
							{Tick: "ADA", Src: enums.Bitkub},
						},
					},
					TrackTimes: 69,
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/track binance btc ada 69",
				TargetBot: enums.TrackBot,
			},
			Expected: BotCommand{
				Track: trackCommand{
					quoteCommand: quoteCommand{
						Securities: []usecase.Security{
							{Tick: "BTC", Src: enums.Binance},
							{Tick: "ADA", Src: enums.Binance},
						},
					},
					TrackTimes: 69,
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/track coinbase btc ada 69",
				TargetBot: enums.TrackBot,
			},
			Expected: BotCommand{
				Track: trackCommand{
					quoteCommand: quoteCommand{
						Securities: []usecase.Security{
							{Tick: "BTC", Src: enums.Coinbase},
							{Tick: "ADA", Src: enums.Coinbase},
						},
					},
					TrackTimes: 69,
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/alert gc=f > 0",
				TargetBot: enums.AlertBot,
			},
			Expected: BotCommand{
				Alert: usecase.Alert{
					Security: usecase.Security{
						Tick: "GC=F",
						Src:  enums.Yahoo,
					},
					Condition: enums.Gt,
					QuoteType: enums.Last,
					Target:    0,
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/alert gc=f bid > 0",
				TargetBot: enums.AlertBot,
			},
			Expected: BotCommand{
				Alert: usecase.Alert{
					Security: usecase.Security{
						Tick: "GC=F",
						Src:  enums.Yahoo,
					},
					Condition: enums.Gt,
					QuoteType: enums.Bid,
					Target:    0,
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/alert satang btc bid > 112",
				TargetBot: enums.AlertBot,
			},
			Expected: BotCommand{
				Alert: usecase.Alert{
					Security: usecase.Security{
						Tick: "BTC",
						Src:  enums.Satang,
					},
					Condition: enums.Gt,
					QuoteType: enums.Bid,
					Target:    112,
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/alert bitkub btc < 112",
				TargetBot: enums.AlertBot,
			},
			Expected: BotCommand{
				Alert: usecase.Alert{
					Security: usecase.Security{
						Tick: "BTC",
						Src:  enums.Bitkub,
					},
					Condition: enums.Lt,
					QuoteType: enums.Last,
					Target:    112,
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/alert bitkub btc bid > 112",
				TargetBot: enums.AlertBot,
			},
			Expected: BotCommand{
				Alert: usecase.Alert{
					Security: usecase.Security{
						Tick: "BTC",
						Src:  enums.Bitkub,
					},
					Condition: enums.Gt,
					QuoteType: enums.Bid,
					Target:    112,
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/alert binance btc bid > 112",
				TargetBot: enums.AlertBot,
			},
			Expected: BotCommand{
				Alert: usecase.Alert{
					Security: usecase.Security{
						Tick: "BTC",
						Src:  enums.Binance,
					},
					Condition: enums.Gt,
					QuoteType: enums.Bid,
					Target:    112,
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/alert coinbase btc bid > 112",
				TargetBot: enums.AlertBot,
			},
			Expected: BotCommand{
				Alert: usecase.Alert{
					Security: usecase.Security{
						Tick: "BTC",
						Src:  enums.Coinbase,
					},
					Condition: enums.Gt,
					QuoteType: enums.Bid,
					Target:    112,
				},
			},
		},
		{
			In: UserCommand{
				Text:      "/alert coinbase btc < 112",
				TargetBot: enums.AlertBot,
			},
			Expected: BotCommand{
				Alert: usecase.Alert{
					Security: usecase.Security{
						Tick: "BTC",
						Src:  enums.Coinbase,
					},
					Condition: enums.Lt,
					QuoteType: enums.Last,
					Target:    112,
				},
			},
		},
	}
)
