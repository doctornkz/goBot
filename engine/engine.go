package engine

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/doctornkz/goBot/formatter"
	log "github.com/sirupsen/logrus"
)

const (
	limitRow  int = 20 // Limit of rows in Top list
	limitChar int = 13 // Limit of characters in Username
)

// DigestMessage structure to formatter package
type DigestMessage struct {
	Header       string
	Words        []string
	Wordsauthors []string
	Newusers     []string
	Leftusers    []string
}

// StatusMessage structure to formatter package
type StatusMessage struct {
	Header  string
	TopList string
}

// User structure, mirror of user table
type User struct {
	UserID      int
	UserName    string
	FirstName   string
	LastName    string
	NumMessages int
	Date        int64 // Time of entering chat or leaving
}

func check(e error) {
	if e != nil {
		log.Error(e)
	}
}

func shorter(username string, limit int) string {
	usernameSlice := strings.Split(username, "")
	if len(usernameSlice) < limit {
		return username
	}

	return strings.Join(usernameSlice[:limit], "") + "..."

}

// SetUser - main user updater
func SetUser(db *sql.DB, user *User) {
	tx, err := db.Begin()
	check(err)

	sqlQueryUser := "insert or replace into user (userid, username, firstname, lastname, num_messages, date) values (?, ?, ?, ?, ?, ?)"

	log.Printf("Engine: SQL Insert %s", sqlQueryUser)
	insertUserState, err := tx.Prepare(sqlQueryUser)
	check(err)
	defer insertUserState.Close()
	_, err = insertUserState.Exec(user.UserID, user.UserName, user.FirstName, user.LastName, user.NumMessages, user.Date)
	check(err)

	tx.Commit()
}

// GetUser - main function around user
func GetUser(db *sql.DB, ID int) (user *User) {
	sqlSelectQuery := "select username, firstname, lastname, num_messages, date from user where userid=?"
	query, err := db.Prepare(sqlSelectQuery)
	if err != nil {
		log.Println("Engine: GetUser failed on Prepare")
		log.Error(err)
	}
	defer query.Close()
	var username string
	var firstname string
	var lastname string
	var nummessages int
	var date int64

	err = query.QueryRow(ID).Scan(&username, &firstname, &lastname, &nummessages, &date)
	if err != nil {
		log.Println("Engine: User is not exist")
		return &User{
			UserID:      ID,
			UserName:    "",
			FirstName:   "",
			LastName:    "",
			NumMessages: -1, //  -1:Left or non-existing user, 0:New user, 0+: Active user
			Date:        time.Now().Unix(),
		}
	}
	log.Println("Engine: User exist")
	return &User{
		UserID:      ID,
		UserName:    username,
		FirstName:   firstname,
		LastName:    lastname,
		NumMessages: nummessages,
		Date:        date,
	}

}

// Status  - TOP20 in chat
func Status(db *sql.DB, ID int) string {

	rows, err := db.Query("select username, firstname, userid, num_messages from user order by num_messages desc;")
	check(err)
	defer rows.Close()
	output := ""
	index := 1
	for rows.Next() { // Generating TOP 20
		var username string
		var firstname string
		var count string
		var userid int
		err = rows.Scan(&username, &firstname, &userid, &count)
		if username == "" {
			username = firstname
		}

		username = shorter(username, limitChar)
		check(err)
		log.Println(strconv.Itoa(index), username, count)

		if index <= limitRow {
			output = output + strconv.Itoa(index) + ". " + username + " = " + count + "\r\n"
		}
		if (ID == userid) && (index > limitRow) {
			output = output + "...\r\n" + strconv.Itoa(index) + ". " + username + " = " + count + "\r\n"
		}
		index++
	}
	err = rows.Err()
	check(err)

	message := &StatusMessage{
		Header:  "-=TOP LIST=-",
		TopList: output}
	messageEncoded, _ := json.Marshal(message)
	return formatter.StatusFormatter(messageEncoded)

}

// Digest generator
func Digest(db *sql.DB, historyhour int64) string {
	period := time.Now().Unix() - historyhour*60*60

	// Select active users
	userrows, err := db.Query("select distinct userid from messages where date>=?", period)
	check(err)
	defer userrows.Close()

	wordsauthors := make([]string, 0, 10)
	users := ""
	for userrows.Next() {
		var userid int
		err = userrows.Scan(&userid)
		username := GetUser(db, userid).UserName

		if username == "" {
			username = GetUser(db, userid).FirstName
		}
		if !strings.Contains(users, username) {
			wordsauthors = append(wordsauthors, username) // Slice with autors
		}
		check(err)
	}

	// Select entered and left users
	statusrows, err := db.Query("select distinct userid, num_messages from user where date>=? ", period)
	check(err)
	defer userrows.Close()

	newusers := make([]string, 0, 10)
	leftusers := make([]string, 0, 10)
	for statusrows.Next() {
		var userid int
		var nummessages int
		err = statusrows.Scan(&userid, &nummessages)
		username := GetUser(db, userid).UserName

		if username == "" {
			username = GetUser(db, userid).FirstName
		}

		if nummessages < 0 {
			leftusers = append(leftusers, username)
		} else {
			newusers = append(newusers, username)
		}
		check(err)
	}

	// Select tags
	rows, err := db.Query("select text from messages where date>=?", period)
	check(err)
	defer rows.Close()
	words := make([]string, 0, 20)
	for rows.Next() {
		var messages string
		err = rows.Scan(&messages)
		for _, word := range strings.Split(messages, " ") {
			words = append(words, word)

		}
	}
	check(err)

	// Marshalizing
	message := &DigestMessage{
		Header:       "*-= DIGEST 12H =-*",
		Words:        words,
		Wordsauthors: wordsauthors,
		Newusers:     newusers,
		Leftusers:    leftusers}
	messageEncoded, _ := json.Marshal(message)
	return formatter.DigestFormatter(messageEncoded)
}
