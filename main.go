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

	bot.Debug = true

	log.Printf("Bot poller: Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	for {
		log.Println("Bot poller: Pre-update section")
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			log.Println("Bot poller: ID section")
			ID := update.Message.From.ID
			log.Println("Bot poller: UserName section")
			UserName := update.Message.From.UserName
			log.Println("Bot poller: FirstName section")
			FirstName := update.Message.From.FirstName
			log.Println("Bot poller: LastName section")
			LastName := update.Message.From.LastName
			log.Println("Bot poller: ChatID section")
			ChatID := update.Message.Chat.ID
			log.Println("Bot poller: Text sections")
			Text := update.Message.Text
			log.Printf("Bot poller: ID: %d UserName: %s FirstName: %s LastName: %s", ID, UserName, FirstName, LastName)
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
