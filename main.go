package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/asjustas/goini"
	"github.com/doctornkz/goBot/engine"
	updater "github.com/doctornkz/goBot/updater"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var dbName string
var dbDriver string
var dir string
var config string
var apiKey string

func init() {
	dirPtr := flag.String("dir", "./", "Working directory")
	confPtr := flag.String("c", "settings.ini", "default config file. See settings.ini.example")
	apiKeyPtr := flag.String("apikey", "", "Bot ApiKey. See @BotFather messages for details")
	dbNamePtr := flag.String("dbname", "", "Database of users")
	dbDriverPtr := flag.String("dbdriver", "sqlite3", "Driver DB.")

	flag.Parse()
	dir = *dirPtr
	config = *confPtr
	apiKey = *apiKeyPtr
	dbName = *dbNamePtr
	dbDriver = *dbDriverPtr
	var err error

	conf, err := goini.Load(dir + config)
	if err != nil {
		log.Printf("Bot poller: Config %s not found, go CLI mode", dir+config)
	}

	//dbNameString := dbName
	if dbName == "" {
		dbName = conf.Str("main", "SQLITE_DB")
		if dbName == "" {
			log.Printf("Bot poller: Something wrong with DB name, %s", dbName)
			log.Panic(err)
		}
	}

	//apiString := apiKey
	if apiKey == "" {
		apiKey = conf.Str("main", "ApiKey")
		if apiKey == "" {
			log.Printf("Bot poller: Something wrong with Apikey file, %s", apiKey)
			log.Panic(err)
		}
	}

	log.Printf("Bot poller: DriverDB %s , NameDB %s", dbDriver, dbName)
	db, err = sql.Open(dbDriver, dbName)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Printf("Bot poller: Something wrong with your key, %s", apiKey)
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
			log.Println("Bot poller: Text sections")
			Date := update.Message.Date
			log.Printf("Bot poller: ID: %d UserName: %s FirstName: %s LastName: %s", ID, UserName, FirstName, LastName)
			if update.Message.IsCommand() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				switch update.Message.Command() {
				case "help":
					msg.Text = "type /sayhi or /status."
				case "sayhi":
					msg.Text = "Hi :)"
				case "status":
					msg.Text = engine.Status(db, ID) // TODO Make limit
				default:
					msg.Text = "I don't know that command"
				}
				bot.Send(msg)
			} else {
				log.Printf("Bot poller: [%s] (ID: %d) %d %s", UserName, ID, ChatID, Text)
				updater.Update(db, ID, UserName, FirstName, LastName, Date, Text) // TODO: make username translation
			}
		}
	}

}
