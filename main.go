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

var config = struct {
	db       *sql.DB
	dbName   string
	dbDriver string
	dir      string
	config   string
	apiKey   string
	chatID   int64
}{
	dbName:   "./empty.db",
	dbDriver: "sqlite3",
	dir:      "./",
	config:   "settings.ini",
	apiKey:   "",
	chatID:   0,
}

func init() {
	dirPtr := flag.String("dir", "./", "Working directory")
	confPtr := flag.String("c", "settings.ini", "default config file. See settings.ini.example")
	apiKeyPtr := flag.String("apikey", "", "Bot ApiKey. See @BotFather messages for details")
	chatIDPtr := flag.Int64("chat", 0, "Chat uniq ID")
	dbNamePtr := flag.String("dbname", "", "Database of users")
	dbDriverPtr := flag.String("dbdriver", "sqlite3", "Driver DB.")

	flag.Parse()

	config.dir = *dirPtr
	config.config = *confPtr
	config.apiKey = *apiKeyPtr
	config.dbName = *dbNamePtr
	config.dbDriver = *dbDriverPtr
	config.chatID = *chatIDPtr

	conf, err := goini.Load(config.dir + config.config)
	if err != nil {
		log.Printf("Bot poller: Config %s not found, go CLI mode", config.dir+config.config)
	}

	//dbNameString := dbName
	if config.dbName == "" {
		config.dbName = conf.Str("main", "SQLITE_DB")
		if config.dbName == "" {
			log.Printf("Bot poller: Something wrong with DB name, %s", config.dbName)
			log.Panic(err)
		}
	}

	//apiString := apiKey
	if config.apiKey == "" {
		config.apiKey = conf.Str("main", "ApiKey")
		if config.apiKey == "" {
			log.Printf("Bot poller: Something wrong with Apikey file, %s", config.apiKey)
			log.Panic(err)
		}
	}

	log.Printf("Bot poller: DriverDB %s , NameDB %s", config.dbDriver, config.dbName)
	config.db, err = sql.Open(config.dbDriver, config.dbName)
	if err != nil {
		log.Fatal(err)
	}

	if err = config.db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	bot, err := tgbotapi.NewBotAPI(config.apiKey)
	if err != nil {
		log.Printf("Bot poller: Something wrong with your key, %s", config.apiKey)
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
			log.Println("Bot poller: Text sections")
			Text := update.Message.Text
			log.Println("Bot poller: Text sections")
			Date := update.Message.Date
			log.Printf("Bot poller: ID: %d UserName: %s FirstName: %s LastName: %s", ID, UserName, FirstName, LastName)
			CurrentChatID := update.Message.Chat.ID
			if CurrentChatID != config.chatID {
				log.Printf("Bot poller: Wrong chat %d ID: %d UserName: %s FirstName: %s LastName: %s", CurrentChatID, ID, UserName, FirstName, LastName)
				continue
			}
			if update.Message.IsCommand() && engine.IfUserExist(config.db, ID) {
				msg := tgbotapi.NewMessage(config.chatID, "")
				switch update.Message.Command() {
				case "help":
					msg.Text = "type /sayhi or /status."
				case "sayhi":
					msg.Text = "Hi :)"
				case "status":
					msg.Text = engine.Status(config.db, ID) // TODO Make limit
				default:
					msg.Text = "I don't know that command"
				}
				bot.Send(msg)
			} else {
				log.Printf("Bot poller: [%s] (ID: %d) %d %s", UserName, ID, config.chatID, Text)
				updater.Update(config.db, ID, UserName, FirstName, LastName, Date, Text) // TODO: make username translation
			}

		}
	}

}
