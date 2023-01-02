package telegram

import (
	"fmt"
	"github.com/Krynegal/Librarian.git/pkg/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
)

type Bot struct {
	bot         *tgbotapi.BotAPI
	storage     storage.Storage
	StateKeeper StateKeeper
}

func NewBot(bot *tgbotapi.BotAPI, storage storage.Storage) *Bot {
	return &Bot{
		bot:     bot,
		storage: storage,
		StateKeeper: StateKeeper{
			mu:     sync.RWMutex{},
			states: map[int]string{},
		},
	}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

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
