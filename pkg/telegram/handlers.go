package telegram

import (
	"encoding/json"
	"fmt"
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

/*
Переменные, указывающее, по какому параметру
будет осуществляться поиск в БД
*/
var (
	waitBookTitle      bool
	waitAuthorLastname bool
)

//Основной обработчик диалоговых сообщений
func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	log.Printf("[%s] %s", message.From.UserName, message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)

	if waitBookTitle {
		//msg.Text, _ = b.storage.GetBooksByTitle(msg.Text)
		msg.Text, _ = b.makeReplyMessage(&msg)
		waitBookTitle = false
		msg.ReplyMarkup = makeChoiceKeyboard()
	}

	if waitAuthorLastname {
		//msg.Text, _ = b.storage.GetBooksByAuthor(msg.Text)
		msg.Text, _ = b.makeReplyMessage(&msg)
		waitAuthorLastname = false
		msg.ReplyMarkup = makeChoiceKeyboard()
	}

	if msg.Text == "Поиск" {
		msg.Text = "Как будем искать?"
		msg.ReplyMarkup = makeChoiceKeyboard()
	}

	if msg.Text == "По автору" {
		msg.Text = "Введите фамилию автора"
		keyboard := tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = keyboard
	}

	if msg.Text == "По названию" {
		msg.Text = "Введите название книги"
		keyboard := tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = keyboard
	}

	if msg.Text == "Введите название книги" {
		waitBookTitle = true
	}

	if msg.Text == "Введите фамилию автора" {
		waitAuthorLastname = true
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
func (b *Bot) makeReplyMessage(msg *tgbotapi.MessageConfig) (string, error) {
	if msg.Text == "null" {
		return "К сожалению, ничего не нашел по вашему запросу", nil
	}
	msgFound := tgbotapi.NewMessage(msg.BaseChat.ChatID, "Вот что я нашёл")
	_, err := b.bot.Send(msgFound)
	if err != nil {
		return "", err
	}
	
	resp, err := makeRepresentativeData(msg.Text)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func makeRepresentativeData(msg string) (string, error) {
	//var b storage.BookInfo
	//var sliceRes []storage.BookInfo
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(msg), &result); err != nil {
		return "", err
	}
	fmt.Println("\nmsg: ", msg)
	fmt.Println("\nresult: ", result)

	var resString string
	format := "Книга: %s\nШкаф: %s\nСекция: %v\nПолка: %v\n"
	if len(result) > 1 {
		format += "\n"
	}
	for _, value := range result {
		resString += fmt.Sprintf(format,
			value["name"], value["bookcase"], value["section_number"], value["shelf_number"])
	}

	return resString, nil
}
