package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v3"

	tghandler "github.com/artnoi43/fngobot/lib/bot/handler/telegram"
	"github.com/artnoi43/fngobot/lib/enums"
	"github.com/artnoi43/fngobot/lib/etc/help"
	"github.com/artnoi43/fngobot/lib/parse"
)

func handle(b *tb.Bot, token string) error {

	log.Printf("initialized bot: %s", token)

	// sigChan for receiving OS signals for graceful shutdowns
	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
		syscall.SIGTERM, // kill -SIGTERM XXXX
	)

	// Graceful shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-sigChan
		log.Println("closed poller connection")
	}()

	sendFail := func() {
		log.Println("error sending Telegram message to recipient")
	}

	b.Handle("/help", handleFunc(b, "/help"))
	b.Handle("/quote", handleFunc(b, "/quote"))
	b.Handle("/track", handleFunc(b, "/track"))
	b.Handle("/alert", handleFunc(b, "/alert"))
	b.Handle("/handlers", handleFunc(b, "/handlers"))

	// Welcome/Greeting
	b.Handle("/start", func(c tb.Context) error {
		log.Println(c.Text())
		if _, err := b.Reply(c.Message(), help.LONG); err != nil {
			sendFail()
		}
		if _, err := b.Reply(c.Message(), "Hello!\nWelcome to FnGoBot chat!"); err != nil {
			sendFail()
		}
		return nil
	})

	// Stop a tracking or alerting Telegram handler
	b.Handle("/stop", func(c tb.Context) error {
		senderId := c.Sender().ID
		uuids := strings.Split(c.Text(), " ")[1:]
		for _, uuid := range uuids {
			// Stop is Handlers method
			idx, ok := tghandler.SenderHandlers[senderId].Stop(uuid)
			if ok {
				// Remove is a plain function
				tghandler.Remove(senderId, idx)
			}
		}
		return nil
	})

	/* Catch-all help message for unhandled text */
	b.Handle(tb.OnText, func(c tb.Context) error {
		log.Println(c.Text())
		if _, err := b.Reply(
			c.Message(),
			fmt.Sprintf("wut? %s\nSee /help for help", c.Text()),
		); err != nil {
			sendFail()
		}
		return nil
	})

	go func() {
		log.Println("fngobot started")
		b.Start()
	}()

	wg.Wait()
	log.Println("fngobot exited")
	return nil
}

func handleFunc(
	b *tb.Bot,
	command enums.InputCommand,
) func(c tb.Context) error {
	return func(c tb.Context) error {
		targetBot, exits := enums.BotMap[command]
		if !exits {
			return fmt.Errorf("invalid command")
		}
		cmd, parseError := parse.UserCommand{
			Text:      c.Text(),
			TargetBot: targetBot,
		}.Parse()
		h := tghandler.New(b, c, &cmd, conf.Telegram)
		if parseError != 0 {
			h.HandleParsingError(parseError)
			return fmt.Errorf("parseError: %d", parseError)
		}
		defer h.Done()
		h.Handle(targetBot)

		if targetBot == enums.HelpBot {
			if _, err := b.Reply(c.Message(), cmd.Help.HelpMessage); err != nil {
				return errors.Wrap(err, "failed to send help message")
			}
			return nil
		}
		return nil
	}
}
