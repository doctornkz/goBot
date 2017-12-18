package engine

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"
)

// User structure, mirror of user table
type User struct {
	UserID      int
	UserName    string
	FirstName   string
	LastName    string
	NumMessages int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// SetUser - main user updater
func SetUser(db *sql.DB, user *User) {
	tx, err := db.Begin()
	check(err)
	if user.NumMessages == -1 {

		sqlQueryUser := "insert into user (userid, username, firstname, lastname, num_messages) values (?, ?, ?, ?, ?)"
		log.Printf("Engine: SQL Insert UserName %s", sqlQueryUser)
		insertUserState, err := tx.Prepare(sqlQueryUser)
		check(err)
		defer insertUserState.Close()
		_, err = insertUserState.Exec(user.UserID, user.UserName, user.FirstName, user.LastName, user.NumMessages)
		check(err)

	} else {
		sqlQueryUser := "update user set username=?, firstname=?, lastname=?, num_messages=? where userid=?"
		log.Printf("Engine: SQL Update UserName %s", sqlQueryUser)
		updateUserState, err := tx.Prepare(sqlQueryUser)
		check(err)
		defer updateUserState.Close()
		_, err = updateUserState.Exec(user.UserName, user.FirstName, user.LastName, user.NumMessages, user.UserID)
		check(err)

	}

	tx.Commit()
}

// GetUser - main function around user
func GetUser(db *sql.DB, ID int) (user *User) {
	sqlSelectQuery := "select username, firstname, lastname, num_messages from user where userid=?"
	query, err := db.Prepare(sqlSelectQuery)
	if err != nil {
		log.Println("Engine: GetUser failed on Prepare")
		log.Fatal(err)
	}
	defer query.Close()
	var username string
	var firstname string
	var lastname string
	var nummessages int

	err = query.QueryRow(ID).Scan(&username, &firstname, &lastname, &nummessages)
	if err != nil {
		log.Println("Engine: User is not exist")
		return &User{
			UserID:      ID,
			UserName:    "",
			FirstName:   "",
			LastName:    "",
			NumMessages: -1, // Not existing user
		}
	}
	log.Println("Engine: User exist")
	return &User{
		UserID:      ID,
		UserName:    username,
		FirstName:   firstname,
		LastName:    lastname,
		NumMessages: nummessages,
	}

}

// Status  - TOP20 in chat
func Status(db *sql.DB, ID int) string {

	rows, err := db.Query("select username, firstname, userid, num_messages from user order by num_messages desc;")
	check(err)
	defer rows.Close()

	output := " -= TOP LIST =- \r\n"
	index := 1
	limit := 20
	for rows.Next() { // Generating TOP 20
		var username string
		var firstname string
		var count string
		var userid int
		err = rows.Scan(&username, &firstname, &userid, &count)
		if username == "" {
			username = firstname
		}

		check(err)
		log.Println(strconv.Itoa(index), username, count)
		if index <= limit {
			output = output + strconv.Itoa(index) + ". " + username + " = " + count + "\r\n"
		}
		if (ID == userid) && (index > limit) {
			output = output + "...\r\n" + strconv.Itoa(index) + ". " + username + " = " + count + "\r\n"
		}
		index++
	}
	err = rows.Err()
	check(err)
	return output
}

// Digest generator
func Digest(db *sql.DB, historyhour int64) string {

	period := time.Now().Unix() - historyhour*60*60

	userrows, err := db.Query("select distinct userid from messages where date>=?", period)
	check(err)
	defer userrows.Close()

	users := ""
	for userrows.Next() {
		var userid int
		err = userrows.Scan(&userid)
		username := GetUser(db, userid).UserName

		if username == "" {
			username = GetUser(db, userid).FirstName
		}
		if !strings.Contains(users, username) {
			users += username + ", "
		}
		check(err)
	}
	rows, err := db.Query("select text from messages where date>=?", period)
	check(err)
	defer rows.Close()

	header := " -= DIGEST 12H =- \r\n"
	cleanedMessage := ""
	for rows.Next() {
		var messages string
		err = rows.Scan(&messages)
		for _, Word := range strings.Split(messages, " ") {
			if !strings.Contains(cleanedMessage, Word) {
				cleanedMessage += Word + ", "
			}

		}
	}
	check(err)
	if cleanedMessage != "" {
		return header + cleanedMessage + "( " + users + ")"
	}
	return "Digest is empty "
}
