package main

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
	"tg/pkg/storage/postgresDB"
	"tg/pkg/telegram"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	token := os.Getenv("TOKEN")
	dataSourceName := os.Getenv("DATA_SOURCE_NAME")
	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	botApi.Debug = true

	db, err := initPostgresDB(dataSourceName)
	storage := postgresDB.NewDatabase(db)

	bot := telegram.NewBot(botApi, storage)
	if err := bot.Start(); err != nil {
		log.Fatal(err)
	}
}

func initPostgresDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	return db, nil
}
