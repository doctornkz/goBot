package main

//  Just demonstration

import (
	"database/sql"
	"flag"
	"strconv"
	"time"

	"github.com/asjustas/goini"
	"github.com/doctornkz/goBot/engine"
	"github.com/doctornkz/goBot/updater"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var (
	ada            updater.AdaConfig
	adaFruitEnable bool
	currentVersion string
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
	// Sqlite DB Configuration:
	dbNamePtr := flag.String("dbname", "", "Database of users")
	dbDriverPtr := flag.String("dbdriver", "sqlite3", "Driver DB.")
	// AdaFruit Configuration:
	adafruitHost := flag.String("adahost", "", "AdaFruit host (graphic service)")
	adafruitPort := flag.String("adaport", "", "Port of AdaFruit collector")
	adafruitUser := flag.String("adauser", "", "Username for AdaFruit autorization")
	adafruitToken := flag.String("adatoken", "", "API token for AdaFruit autorization")
	adafruitTopic := flag.String("adatopic", "", "AdaFruit topic")

	flag.Parse()

	config.dir = *dirPtr
	config.config = *confPtr
	config.apiKey = *apiKeyPtr
	config.dbName = *dbNamePtr
	config.dbDriver = *dbDriverPtr
	config.chatID = *chatIDPtr
	// AdaFruit parameters initialize :

	ada.AdafruitHost = *adafruitHost
	ada.AdafruitPort = *adafruitPort
	ada.AdafruitUser = *adafruitUser
	ada.AdafruitToken = *adafruitToken
	ada.AdaFruitTopic = *adafruitTopic

	// TODO: Use switch/select Luke! :
	if ada.AdafruitHost == "" || ada.AdafruitPort == "" || ada.AdafruitUser == "" || ada.AdafruitToken == "" || ada.AdaFruitTopic == "" {
		log.Printf("Configuration: Config for AdaFruit not found, skipping...")
		adaFruitEnable = false
	} else {
		log.Printf("Configuration: Config for AdaFruit found, username %s setting...", ada.AdafruitUser)
		adaFruitEnable = true
	}

	conf, err := goini.Load(config.dir + config.config)
	if err != nil {
		log.Printf("Bot poller: Config %s not found, go CLI mode", config.dir+config.config)
	}

	if config.dbName == "" {
		config.dbName = conf.Str("main", "SQLITE_DB")
		if config.dbName == "" {
			log.Printf("Bot poller: Something wrong with DB name, %s", config.dbName)
			log.Panic(err)
		}
	}

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
	// AdaFruit channel init
	adaChanMessage := make(chan string)
	go updater.ExportStats(ada, adaChanMessage)

	// Telegram init
	bot, err := tgbotapi.NewBotAPI(config.apiKey)
	if err != nil {
		log.Printf("Bot poller: Something wrong with your key, %s", config.apiKey)
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Bot poller: Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	var chatConfig tgbotapi.ChatConfig
	chatConfig.ChatID = config.chatID
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
			newusers := update.Message.NewChatMembers
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
				case "version":
					msg.Text = "GoBot," + currentVersion + ", gh:doctornkz:goBot."
				case "digest12h":
					// msg.ParseMode = "Markdown"              // Markdown works only for Digest, may be bug
					// Markdown removed, generates error:
					// Bad Request: can't parse entities: Can't find end of the entity starting at byte offset 695
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
					log.Println("Bot Poller: Users left from Chat")
					if adaFruitEnable {
						count, err := bot.GetChatMembersCount(chatConfig)
						check(err)
						adaChanMessage <- strconv.Itoa(count)
					}

					user := engine.GetUser(config.db, leftuser.ID)
					user.UserID = leftuser.ID
					user.UserName = leftuser.UserName
					user.FirstName = leftuser.FirstName
					user.LastName = leftuser.LastName
					user.Date = time.Now().Unix()
					user.NumMessages = -1
					engine.SetUser(config.db, user)
					log.Println("Bot Poller: User " + leftuser.FirstName + " go out from Chat")

				} else if newusers != nil {
					log.Println("Bot Poller: Users entered in Chat")
					if adaFruitEnable {
						count, err := bot.GetChatMembersCount(chatConfig)
						check(err)
						adaChanMessage <- strconv.Itoa(count)
					}

					for _, newuser := range *newusers {
						user := engine.GetUser(config.db, newuser.ID)
						user.UserID = newuser.ID
						user.UserName = newuser.UserName
						user.FirstName = newuser.FirstName
						user.LastName = newuser.LastName
						user.Date = time.Now().Unix() // FIXME Double check time with non-existing user
						user.NumMessages = 0
						engine.SetUser(config.db, user)
						log.Println("Bot Poller: User " + user.FirstName + "entered in Chat")
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
