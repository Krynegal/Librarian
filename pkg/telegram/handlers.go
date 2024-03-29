package telegram

import (
	"fmt"
	"github.com/Krynegal/Librarian.git/pkg/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const commandStart = "start"

//Обработчик "команд" - сообщений формата /<any text>
func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

//Основной обработчик диалоговых сообщений
func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	log.Printf("[%s] %s", message.From.UserName, message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)

	state := b.StateKeeper.state(int(message.From.ID))
	log.Printf("\nstate: %s\n", state)

	switch state {
	case "search":
		msg.Text = "Как будем искать?"
		msg.ReplyMarkup = makeChoiceKeyboard()
		log.Printf("\nmsg.Text: %s\n", msg.Text)
		if msg.Text == "По автору" {
			b.StateKeeper.update(int(message.From.ID), "byAuthor")
		} else {
			b.StateKeeper.update(int(message.From.ID), "byTitle")
		}
	case "byAuthor":
		msg.Text = "Введите фамилию автора"
		keyboard := tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = keyboard
		b.StateKeeper.update(int(message.From.ID), "enterAuthorLastname")
	case "byTitle":
		msg.Text = "Введите название книги"
		keyboard := tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = keyboard
		b.StateKeeper.update(int(message.From.ID), "enterBookTitle")
	case "enterBookTitle":
		bookSlice, _ := b.storage.GetBooksByTitle(msg.Text)
		msg.Text, _ = b.makeReplyMessage(&msg, bookSlice)
		msg.ReplyMarkup = makeChoiceKeyboard()
	case "enterAuthorLastname":
		bookSlice, _ := b.storage.GetBooksByAuthor(msg.Text)
		msg.Text, _ = b.makeReplyMessage(&msg, bookSlice)
		msg.ReplyMarkup = makeChoiceKeyboard()
	}

	if msg.Text != message.Text {
		_, err := b.bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
Обработчик команды /start
Отправляет кнопку "Поиск" при запуске бота
*/
func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Чтобы приступить к поиску, нажмите кнопку Поиск")
	findButton := tgbotapi.NewKeyboardButton("Поиск")

	keyboard := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{findButton})
	msg.ReplyMarkup = keyboard

	_, err := b.bot.Send(msg)
	b.StateKeeper.update(int(message.From.ID), "search")
	return err
}

// обработчик неизвестных команд
func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды")
	_, err := b.bot.Send(msg)
	return err
}

//Возвращает объект клавиатуры с двумя кнопками
func makeChoiceKeyboard() tgbotapi.ReplyKeyboardMarkup {
	buttons := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("По автору"),
		tgbotapi.NewKeyboardButton("По названию"),
	}
	keyboard := tgbotapi.NewReplyKeyboard(buttons)
	return keyboard
}

/*
Если результат запроса к БД НЕ null, отправить сообщение об успешном поиске
и показать результат запроса.
Иначе - отправить сообщение о неуспешном поиске
*/
func (b *Bot) makeReplyMessage(msg *tgbotapi.MessageConfig, books []storage.BookInfo) (string, error) {
	if msg.Text == "null" {
		return "К сожалению, ничего не нашел по вашему запросу", nil
	}
	msgFound := tgbotapi.NewMessage(msg.BaseChat.ChatID, "Вот что я нашёл")
	_, err := b.bot.Send(msgFound)
	if err != nil {
		return "", err
	}

	resp := printBookInfo(books)
	return resp, nil
}

func printBookInfo(books []storage.BookInfo) string {
	var resString string
	format := "Книга: %s\nШкаф: %s\nСекция: %v\nПолка: %v\n"
	for i, b := range books {
		if i != len(books) {
			format += "\n"
		}
		resString += fmt.Sprintf(format, b.Name, b.Bookcase, b.SectionNumber, b.ShelfNumber)
	}
	return resString
}
