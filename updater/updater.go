package updater

import (
	"database/sql"
	"log"
	"regexp"
	"strings"

	"github.com/doctornkz/goBot/engine"
	// Snowball stemmer
	"github.com/kljensen/snowball"
	// SQLite3 package
	_ "github.com/mattn/go-sqlite3"
)

// Update message's counter

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Update - update function about user activity
func Update(db *sql.DB, ID int, UserName string, FirstName string, LastName string, Date int, Text string) {
	// Select rows with ID
	user := engine.GetUser(db, ID)

	if user.NumMessages == -1 {
		user.UserID = ID
		user.UserName = UserName
		user.FirstName = FirstName
		user.LastName = LastName
		user.NumMessages = -1
		engine.SetUser(db, user)
		messagesUpdate(db, ID, Date, Text)

	} else {

		messagesUpdate(db, ID, Date, Text)

	}

	currentNumMessages := user.NumMessages
	log.Printf("Updater: %s, %d", user.UserName, currentNumMessages)
	user.NumMessages = currentNumMessages + 1
	log.Printf("Updater: %s, %d", user.UserName, user.NumMessages)
	log.Println(user)

	engine.SetUser(db, user)

}

func checkWord(db *sql.DB, Word string) int {
	query, err := db.Prepare("select categoryid from words where word = ?")
	check(err)
	defer query.Close()
	var category int
	err = query.QueryRow(Word).Scan(&category)
	if err != nil {
		return -10
	}
	return category
}

func insertWord(db *sql.DB, Word string, Date int) {
	tx, err := db.Begin()
	check(err)

	query, err := tx.Prepare("insert into words(word, categoryid, userid) values(?, ?, ?)")

	check(err)
	defer query.Close()

	_, err = query.Exec(Word, -1, 0)
	check(err)
	tx.Commit()
}

// MessagesUpdate - updater messages in chat
func messagesUpdate(db *sql.DB, ID int, Date int, Text string) {
	log.Println("Updater: Messages insert section")
	tx, err := db.Begin()
	check(err)

	cleanedMessage := ""
	for _, Word := range strings.Split(Text, " ") {
		if !strings.Contains(cleanedMessage, Word) {
			if checkWord(db, stemming(Word)) == 1 { // FIXME: Hardcoded category
				var re = regexp.MustCompile(`[a-z]|[@$%&*~#=/_"!?. ,:;\-\\+1234567890(){}\[\]]`)
				Word = re.ReplaceAllString(Word, "")

				sqlMessQuery := "insert into messages (userid, date, text) values (?, ?, ?)"
				log.Printf("Updater: SQL Insert %s", sqlMessQuery)

				updateMessageState, err := tx.Prepare(sqlMessQuery)
				check(err)
				defer updateMessageState.Close()

				_, err = updateMessageState.Exec(ID, Date, Word)
				check(err)

				log.Println("Updater: Message committed")
			}
		}
	}
	tx.Commit()
	log.Printf("Updater: Phrase to save: %s", cleanedMessage)

}

func stemming(Word string) string {
	var re = regexp.MustCompile(`[a-z]|[@$%&*~#=/_"!?. ,:;\-\\+1234567890(){}\[\]]`)
	Word = re.ReplaceAllString(Word, "")
	stemmed, err := snowball.Stem(Word, "russian", true)
	check(err)
	return stemmed
}
