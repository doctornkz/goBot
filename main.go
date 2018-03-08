package main

import (
	"database/sql"
	"flag"
	"time"

	"github.com/asjustas/goini"
	"github.com/doctornkz/goBot/engine"
	"github.com/doctornkz/goBot/updater"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var config = struct {
	db       *sql.DB
	dbName   string
	dbDriver string
	dir      string
	config   string
	apiKey   string
	chatID   int64
	logfile  string
}{
	dbName:   "./empty.db",
	dbDriver: "sqlite3",
	dir:      "./",
	config:   "settings.ini",
	apiKey:   "",
	chatID:   0,
	logfile:  "./gobot.log",
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

func check(e error) {
	if e != nil {
		log.Error(e)
	}
}

func main() {

	bot, err := tgbotapi.NewBotAPI(config.apiKey)
	if err != nil {
		log.Printf("Bot poller: Something wrong with your key, %s", config.apiKey)
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Bot poller: Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	check(err)
	for {
		log.Println("Bot poller: Pre-update section")
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			ID := update.Message.From.ID
			user := engine.GetUser(config.db, ID)
			username := update.Message.From.UserName
			firstname := update.Message.From.FirstName
			lastname := update.Message.From.LastName
			text := update.Message.Text
			date := update.Message.Date
			newuser := update.Message.NewChatMembers
			leftuser := update.Message.LeftChatMember
			log.Printf("Bot poller: ID: %d UserName: %s FirstName: %s LastName: %s", ID, username, firstname, lastname)
			currentChatID := update.Message.Chat.ID
			if currentChatID != config.chatID {
				log.Printf("Bot poller: Wrong chat %d ID: %d UserName: %s FirstName: %s LastName: %s", currentChatID, ID, username, firstname, lastname)
				continue
			}
			// Command messages
			if update.Message.IsCommand() && user.NumMessages > 0 {
				msg := tgbotapi.NewMessage(config.chatID, "")

				switch update.Message.Command() {
				case "help":
					msg.Text = "type /sayhi, /digest12h or /status."
				case "sayhi":
					msg.Text = "Hi :)"
				case "digest12h":
					msg.ParseMode = "Markdown"              // Markdown works only for Digest, may be bug
					msg.Text = engine.Digest(config.db, 12) // To do hours and ID
				case "status":
					msg.Text = engine.Status(config.db, ID) // Make limit (1..20, all)
				default:
					continue
				}
				_, err := bot.Send(msg)
				check(err)

			} else {
				log.Printf("Bot poller: [%s] (ID: %d) %d %s", username, ID, config.chatID, text)
				// Check new and left users:  // TODO: Lots the duplicates, replace to function?
				if leftuser != nil {
					user := engine.GetUser(config.db, leftuser.ID)
					user.UserID = leftuser.ID
					user.UserName = leftuser.UserName
					user.FirstName = leftuser.FirstName
					user.LastName = leftuser.LastName
					user.Date = time.Now().Unix()
					user.NumMessages = -1
					engine.SetUser(config.db, user)
					log.Println("Bot Poller: User " + leftuser.FirstName + " go out from Chat")

				} else if newuser != nil {
					for _, newuservalue := range *newuser {
						user := engine.GetUser(config.db, newuservalue.ID)
						user.UserID = newuservalue.ID
						user.UserName = newuservalue.UserName
						user.FirstName = newuservalue.FirstName
						user.LastName = newuservalue.LastName
						user.Date = time.Now().Unix() // FIXME Double check time with non-existing user
						user.NumMessages = 0
						engine.SetUser(config.db, user)
						log.Println("Bot Poller: User " + newuservalue.FirstName + " entered in Chat")
					}

				} else {
					// Update chat trivial messaging
					log.Println("Bot Poller: Start updating. Updater come in.")
					updater.Update(config.db, ID, username, firstname, lastname, int64(date), text) // TODO: make username translation
				}
			}

		}
	}

}
