package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"tg/pkg/storage"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	storage storage.Storage
}

func NewBot(bot *tgbotapi.BotAPI, storage storage.Storage) *Bot {
	return &Bot{
		bot: bot,
		storage: storage,
	}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		//Handle commands
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				fmt.Println(err)
			}
			continue
		}

		//Handle other messages
		if err := b.handleMessage(update.Message); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
