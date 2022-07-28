package main

import (
	"database/sql"
	"github.com/Krynegal/Librarian.git/pkg/storage/postgresDB"
	"github.com/Krynegal/Librarian.git/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
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
