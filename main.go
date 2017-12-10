package main

import (
	"log"

	"github.com/doctornkz/goBot/engine"

	_ "github.com/mattn/go-sqlite3"

	"github.com/asjustas/goini"
	updater "github.com/doctornkz/goBot/updater"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	conf, err := goini.Load("./settings.ini")
	if err != nil {
		panic(err)
	}

	apiString := conf.Str("main", "ApiKey")
	bot, err := tgbotapi.NewBotAPI(apiString)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Bot poller: Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	for {
		select {
		case update := <-updates:
			ID := update.Message.From.ID
			UserName := update.Message.From.UserName
			FirstName := update.Message.From.FirstName
			LastName := update.Message.From.LastName

			ChatID := update.Message.Chat.ID
			Text := update.Message.Text

			if update.Message.IsCommand() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				switch update.Message.Command() {
				case "help":
					msg.Text = "type /sayhi or /status."
				case "sayhi":
					msg.Text = "Hi :)"
				case "status":
					msg.Text = engine.Status(ID) // TODO Make limit
				default:
					msg.Text = "I don't know that command"
				}
				bot.Send(msg)
			} else {
				log.Printf("Bot poller: [%s] (ID: %d) %d %s", UserName, ID, ChatID, Text)
				updater.Update(ID, UserName, FirstName, LastName) // TODO: make username translation
			}
		}
	}

}
