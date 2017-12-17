package engine

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"
	//logger "github.com/doctornkz/goBot/logger"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// IfUserExist check user valid in chat
func IfUserExist(db *sql.DB, ID int) bool {
	// Select rows with ID
	sqlSelectQuery := "select count from num_messages where userid= ?"
	query, err := db.Prepare(sqlSelectQuery)
	if err != nil {
		log.Println("Engine: ifUserExist false on Prepare")
		log.Fatal(err)
		return false
	}
	defer query.Close()
	// Query section
	var count int
	err = query.QueryRow(ID).Scan(&count)
	if err != nil {
		log.Println("Engine: ifUserExist false on Scan")
		log.Println(err)
		return false
	}
	log.Println("Engine: ifUserExist true")
	return true

}

// Status  - TOP20 in chat
func Status(db *sql.DB, ID int) string {
	rows, err := db.Query("select user.username, user.firstname, num_messages.userid, num_messages.count from user inner join num_messages on user.userid=num_messages.userid order by count desc")
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

	now := time.Now().Unix()
	period := now - historyhour*60*60
	log.Println(period)

	///Date := time
	rows, err := db.Query("select userid, text from messages where date>=?", period)
	check(err)
	defer rows.Close()

	header := " -= DIGEST 12H =- \r\n"
	cleanedMessage := ""
	flooders := ""
	for rows.Next() {
		var messages string
		var userid int
		err = rows.Scan(&userid, &messages)
		for _, Word := range strings.Split(messages, " ") {
			username := userid2Name(db, userid)
			if !strings.Contains(cleanedMessage, Word) {
				cleanedMessage += Word + " "
			}
			if !strings.Contains(flooders, username) {
				flooders += username + " "
			}
		}
	}
	check(err)
	if cleanedMessage != "" {
		return header + cleanedMessage + "( " + flooders + ")"
	}
	return "Digest is empty "
}

func userid2Name(db *sql.DB, ID int) string {
	// Select rows with ID
	sqlSelectQuery := "select username, firstname from user where userid=?"
	query, err := db.Prepare(sqlSelectQuery)
	check(err)
	defer query.Close()

	// Query section
	var username string
	var firstname string
	err = query.QueryRow(ID).Scan(&username, &firstname)
	check(err)
	if username != "" {
		return username
	}

	return firstname

}
