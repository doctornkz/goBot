package updater

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

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

func Update(db *sql.DB, ID int, UserName string, FirstName string, LastName string, Date int, Text string) {
	// Select rows with ID
	sqlSelectQuery := "select count from num_messages where userid= ?"
	log.Printf("Updater: SQL Select %s", sqlSelectQuery)
	stmt, err := db.Prepare(sqlSelectQuery)
	check(err)
	defer stmt.Close()

	// Query section
	//var userid int
	var count int
	operation := "update" // TODO Make simplier, or goto?
	err = stmt.QueryRow(ID).Scan(&count)
	if err != nil {
		// New user detected
		fmt.Println(err)
		//userid = ID
		count = 0
		operation = "insert"
		userUpdate(db, ID, UserName, FirstName, LastName)
	}

	// Update messages
	log.Printf("Updater: Place new phrase in DB")

	messagesUpdate(db, ID, Date, Text)

	log.Printf("Updater: UserID %d found with %d messages", ID, count)

	// Insert section
	log.Printf("Updater: Operation flag is %s", operation)
	// Begin section
	tx, err := db.Begin()
	check(err)

	var sqlQuery string
	// Prepare section

	if operation == "insert" {

		sqlQuery = "insert into num_messages (userid, count) values (?, ?)"
		log.Printf("Updater: SQL Insert %s", sqlQuery)

	} else {

		sqlQuery = "update num_messages set count =? where userid = ?"
		log.Printf("Updater: SQL Update %s", sqlQuery)
	}

	smth, err := tx.Prepare(sqlQuery)
	check(err)
	defer smth.Close()
	// Exec section

	count++
	log.Printf("Updater: UserID: %d count: %d", ID, count)
	if operation == "insert" {
		_, err = smth.Exec(ID, count)
		check(err)
	} else {
		_, err = smth.Exec(count, ID)
		check(err)
	}

	// Commit section
	log.Println("Updater: Pre commit step")
	tx.Commit()
	log.Println("Updater: Committed")
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

	//
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

// UserUpdate  Username updating
func userUpdate(db *sql.DB, ID int, UserName string, FirstName string, LastName string) {

	tx, err := db.Begin()
	check(err)

	sqlQueryUser := "insert into user (userid, username, firstname, lastname) values (?, ?, ?, ?)"
	log.Printf("Updater: SQL Insert UserName %s", sqlQueryUser)

	updateUserState, err := tx.Prepare(sqlQueryUser)
	check(err)
	defer updateUserState.Close()

	_, err = updateUserState.Exec(ID, UserName, FirstName, LastName)
	check(err)
	tx.Commit()
}

func stemming(Word string) string {
	var re = regexp.MustCompile(`[a-z]|[@$%&*~#=/_"!?. ,:;\-\\+1234567890(){}\[\]]`)
	Word = re.ReplaceAllString(Word, "")
	stemmed, err := snowball.Stem(Word, "russian", true)
	check(err)
	return stemmed
}
